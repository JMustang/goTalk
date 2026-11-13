package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var messageReq MessageReq
	err := decoder.Decode(&messageReq)
	if err != nil {
		response := MessageRes{
			Status:  "error",
			Message: "❌ can not decode message",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if messageReq.Message == "" {
		response := MessageRes{
			Status:  "error",
			Message: "❌ message empty",
		}
		responseByte, err := json.Marshal(response)
		if err != nil {
			log.Println("❌ can not marshal response")
			http.Error(w, "❌ can not marshal response", http.StatusInternalServerError)
			return
		}
		http.Error(w, string(responseByte), http.StatusBadRequest)
		return
	}
	response := MessageRes{
		Status:  "success",
		Message: messageReq.Message,
	}
	json.NewEncoder(w).Encode(response)

}

func initSocket() {
	origin := "http://localhost/"
	url := "ws://localhost:3000/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s. \n", msg[:n])
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.HandleFunc("/api/send", handler)

	initSocket()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
