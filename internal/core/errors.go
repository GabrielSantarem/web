package core

import "errors"

var (
	ErrUserNotFound      = errors.New("user não encontrado")
	ErrEmailAlreadyUsed  = errors.New("email já cadastrado")
	ErrUserAlreadyActive = errors.New("user já está ativo")
	ErrInvalidInput      = errors.New("entrada inválida")
	ErrUserAlreadyExists = errors.New("user já existe")
)
