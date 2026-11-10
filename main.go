package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MessageReq struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type MessageRes struct {
	Status  string
	Message string
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var messageReq MessageReq
	err := decoder.Decode(&messageReq)
	if err != nil {
		response := MessageRes{
			Status:  "error",
			Message: "Not right format",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := MessageRes{
			Status:  "success",
			Message: messageReq.Message,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// http.HandleFunc("/hello", handler)
	http.HandleFunc("/api/send", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
