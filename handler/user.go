package handler

import "net/http"

func HandleCreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create User"))
	})
}
