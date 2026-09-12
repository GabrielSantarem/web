package handler

import "github.com/go-playground/validator/v10"

func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	// Verifica se o erro é de fato uma falha de validação do go-playground/validator
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			// fieldErr.Field() traz o nome do campo da Struct (ex: "Name")
			// fieldErr.Tag() traz a regra que quebrou (ex: "min", "email")
			switch fieldErr.Tag() {
			case "required":
				errors[fieldErr.Field()] = "este campo é obrigatório"
			case "min":
				errors[fieldErr.Field()] = "o valor fornecido é menor que o mínimo exigido"
			case "email":
				errors[fieldErr.Field()] = "o formato do e-mail é inválido"
			default:
				errors[fieldErr.Field()] = "validação falhou na regra: " + fieldErr.Tag()
			}
		}
		return errors
	}

	// Caso ocorra outro tipo de erro de validação inesperado
	return map[string]string{"error": err.Error()}
}
