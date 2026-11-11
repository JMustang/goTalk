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

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// http.HandleFunc("/hello", handler)
	http.HandleFunc("/api/send", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
