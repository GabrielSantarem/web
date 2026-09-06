package main

import (
	"net/http"

	"codeberg.org/MrTomate/web/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.HandleHello())
	mux.HandleFunc("/me/{Name}", handler.HandleMe())
	mux.HandleFunc("/", handler.HandleNotFound())

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	})

	http.ListenAndServe(":8080", mux)
}
