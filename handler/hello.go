package handler

import (
	"encoding/json"
	"net/http"
)

func HandleHello() http.HandlerFunc {

	s := Message{Message: "Hello", Id: 1}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(s)
	}
}

func HandleMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		name := r.PathValue("Name")

		me := Me{
			Name: name,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(me)
	}
}

func HandleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
