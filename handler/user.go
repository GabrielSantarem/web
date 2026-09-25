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
			if errors.Is(err, core.ErrUserAlreadyExists) || errors.Is(err, core.ErrEmailAlreadyUsed) {
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

func HandleGetUser(srv core.UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "id obrigatorio",
			})
			return
		}

		user, err := srv.GetUser(id)
		if err != nil {
			if errors.Is(err, core.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{
					Message: "usuario nao encontrado",
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	})
}

func HandleActivateUser(srv core.UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "id obrigatorio",
			})
			return
		}

		err := srv.ActivateUser(id)
		if err != nil {
			if errors.Is(err, core.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{
					Message: "usuario nao encontrado",
					Err:     err.Error(),
				})
				return
			}

			if errors.Is(err, core.ErrUserAlreadyActive) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(ErrorResponse{
					Message: "usuario ja esta ativo",
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "usuario ativado com sucesso",
		})
	})
}

