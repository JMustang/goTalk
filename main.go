package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

type MessageReq struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type MessageRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
	inbox   chan []byte
	mu      sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		inbox:   make(chan []byte),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.join:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.leave:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.inbox:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var messageReq MessageReq
	if err := decoder.Decode(&messageReq); err != nil {
		response := MessageRes{Status: "error", Message: "❌ can not decode message"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if strings.TrimSpace(messageReq.Message) == "" {
		response := MessageRes{Status: "error", Message: "❌ message empty"}
		responseByte, err := json.Marshal(response)
		if err != nil {
			log.Println("❌ can not marshal response")
			http.Error(w, "❌ can not marshal response", http.StatusInternalServerError)
			return
		}
		http.Error(w, string(responseByte), http.StatusBadRequest)
		return
	}

	response := MessageRes{Status: "success", Message: messageReq.Message}
	json.NewEncoder(w).Encode(response)
}

func socketHandler(hub *Hub) websocket.Handler {
	return func(ws *websocket.Conn) {
		client := &Client{conn: ws, send: make(chan []byte, 256)}
		hub.join <- client
		defer func() {
			hub.leave <- client
			ws.Close()
		}()

		go func() {
			for message := range client.send {
				if err := websocket.Message.Send(ws, string(message)); err != nil {
					log.Printf("write websocket message: %v", err)
					return
				}
			}
		}()

		for {
			var payload MessageReq
			if err := websocket.JSON.Receive(ws, &payload); err != nil {
				log.Printf("read websocket message: %v", err)
				return
			}

			if strings.TrimSpace(payload.Message) == "" {
				continue
			}

			response, err := json.Marshal(MessageRes{Status: "success", Message: payload.Message})
			if err != nil {
				log.Printf("marshal message: %v", err)
				continue
			}

			hub.inbox <- response
		}
	}
}

func initSocket() {
	if os.Getenv("GO_TALK_ENABLE_WS") != "true" {
		return
	}

	origin := "http://localhost/"
	url := "ws://localhost:3000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Printf("websocket disabled: %v", err)
		return
	}
	defer ws.Close()

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Printf("failed to write websocket message: %v", err)
		return
	}

	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Printf("failed to read websocket message: %v", err)
		return
	}

	log.Printf("Received: %s", msg[:n])
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	hub := newHub()
	go hub.run()

	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/chat.html", http.StatusMovedPermanently)
	})
	http.HandleFunc("/api/send", handler)
	http.Handle("/ws", socketHandler(hub))

	initSocket()
	log.Println("goTalk running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
