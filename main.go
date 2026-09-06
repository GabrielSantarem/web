package main

import (
	"encoding/json"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloHandler())

	http.ListenAndServe(":8080", mux)
}

type Message struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

func HelloHandler() http.HandlerFunc {
	s := Message{
		Id:      1,
		Message: "Hello world",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(s)
	}
}
