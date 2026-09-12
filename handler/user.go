package handler

import (
	"encoding/json"
	"net/http"

	"codeberg.org/MrTomate/web/internal/core"
	"github.com/go-playground/validator/v10"
)

var validatorInstance = validator.New()

func HandleCreateUser(srv core.UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newUser UserCreateRequest

		err := json.NewDecoder(r.Body).Decode(&newUser)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		if err = validatorInstance.Struct(newUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		user, err := srv.CreateUser(newUser.ToUser())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
}
