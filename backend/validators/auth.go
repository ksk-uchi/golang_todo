package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return map[string]string{"error": "An unexpected error occurred during validation."}
		}

		errorMessages := make(map[string]string)
		for _, fe := range validationErrors {
			field := strings.ToLower(fe.Field())
			switch field {
			case "email":
				switch fe.Tag() {
				case "required":
					errorMessages[field] = "メールアドレスは必須です"
				case "email":
					errorMessages[field] = "有効なメールアドレスを入力してください"
				}
			case "password":
				switch fe.Tag() {
				case "required":
					errorMessages[field] = "パスワードは必須です"
				}
			}
		}
		return errorMessages
	}
	return nil
}

type LoginResponse struct {
	Message string `json:"message"`
}
