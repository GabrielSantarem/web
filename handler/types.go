package handler

import (
	"errors"

	"codeberg.org/MrTomate/web/internal/core"
)

var (
	ErrValidationFailure = errors.New("validation failure")
)

// Message representa a resposta simples do endpoint hello.
type Message struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

// Me representa a resposta com o nome extraído da rota.
type Me struct {
	Name string `json:"name"`
}

type UserCreateRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
}

func (r *UserCreateRequest) ToUser() *core.User {
	return &core.User{
		Name:  r.Name,
		Email: r.Email,
	}
}

// ErrorResponse padroniza todas as respostas de erro da API.
type ErrorResponse struct {
	Message string            `json:"message"`
	Err     string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.Err != "" {
		return e.Err
	}
	return e.Message
}
