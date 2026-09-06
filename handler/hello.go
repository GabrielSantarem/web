package handler

import (
	"encoding/json"
	"net/http"

	"codeberg.org/MrTomate/web/types"
)

func HelloHandler() http.HandlerFunc {

	s := types.Message{Message: "Hello", Id: 1}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(s)
	}
}
