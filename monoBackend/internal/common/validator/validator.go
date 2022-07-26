package validator

import (
	"github.com/go-playground/locales/en"
	uniTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"strings"
)

type ValidationError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func ValidateStruct(input interface{}) []*ValidationError {
	validate := validator.New()
	engl := en.New()
	uniTranslator := uniTrans.New(engl, engl)

	translator, _ := uniTranslator.GetTranslator("en")

	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil
	}

	if err := validate.Struct(input); err != nil {
		return buildTranslatedErrorMessages(err.(validator.ValidationErrors), translator)
	}
	return nil
}

func buildTranslatedErrorMessages(err validator.ValidationErrors, translator uniTrans.Translator) []*ValidationError {
	var errors []*ValidationError

	for _, err := range err {
		var element ValidationError
		element.Key = strings.ToLower(err.Field())
		element.Message = err.Translate(translator)
		errors = append(errors, &element)
	}

	return errors
}
