package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=200"`
}

func (r *CreateTodoRequest) Validate() map[string]string {
	if err := validate.Struct(r); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			// A non-validation error occurred, which is unexpected.
			// Depending on the desired error handling strategy,
			// you might want to log this or return a generic error message.
			return map[string]string{"error": "An unexpected error occurred during validation."}
		}

		errorMessages := make(map[string]string)
		for _, fe := range validationErrors {
			field := strings.ToLower(fe.Field())
			switch field {
			case "title":
				switch fe.Tag() {
				case "required":
					errorMessages[field] = "タイトルは必須です"
				case "max":
					errorMessages[field] = "タイトルは" + fe.Param() + "文字以内で入力してください"
				}
			case "description":
				switch fe.Tag() {
				case "max":
					errorMessages[field] = "説明は" + fe.Param() + "文字以内で入力してください"
				}
			}
		}
		return errorMessages
	}
	return nil
}
