package handler

import "codeberg.org/MrTomate/web/internal/core"

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
