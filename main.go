package main

import (
	"net/http"

	"codeberg.org/MrTomate/web/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HelloHandler())

	http.ListenAndServe(":8080", mux)
}
