package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"codeberg.org/MrTomate/web/internal/core"
)

func HandleCreateUser(srv core.UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newUser UserCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "json invalido",
				Err:     err.Error(),
			})
			return
		}

		if err := validatorInstance.Struct(newUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: ErrValidationFailure.Error(),
				Err:     ErrValidationFailure.Error(),
				Errors:  formatValidationError(err),
			})
			return
		}

		user, err := srv.CreateUser(newUser.ToUser())
		if err != nil {
			if errors.Is(err, core.ErrUserAlreadyExists) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(ErrorResponse{
					Message: "email ja cadastrado",
					Err:     err.Error(),
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "erro interno no servidor",
				Err:     err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
}
