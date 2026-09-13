package handler

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validatorInstance = validator.New()

func formatValidationError(err error) map[string]string {
	errMap := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := strings.ToLower(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				errMap[field] = "campo é obrigatório"
			case "min":
				errMap[field] = "valor fornecido é menor que o mínimo exigido"
			case "email":
				errMap[field] = "formato do e-mail é inválido"
			case "max":
				errMap[field] = "formato acima do máximo permitido"
			default:
				errMap[field] = "validação falhou: " + fieldErr.Tag()
			}
		}
		return errMap
	}

	return map[string]string{"error": err.Error()}
}
