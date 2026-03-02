package validators

import (
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	"reflect"
	"strings"
)

var (
	validate   = validator.New()
	translator ut.Translator
)

func init() {
	japanese := ja.New()
	uni := ut.New(japanese, japanese)
	translator, _ = uni.GetTranslator("ja")
	_ = ja_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag as field name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func TranslateError(err error) map[string]string {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "An unexpected error occurred during validation."}
	}

	errorMessages := make(map[string]string)
	for _, fe := range validationErrors {
		errorMessages[fe.Field()] = fe.Translate(translator)
	}
	return errorMessages
}
