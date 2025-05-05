package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Error string `json:"error"`
}

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

const (
	minNameLength   = 2
	maxNameLength   = 100
	minIDValue      = 1
	nameRegexString = `^[a-zA-Z0-9\s\-_.,!?()]{2,100}$`
)

func New() (*Validator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()

	// Register json tag names
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	if err := registerCustomValidations(validate); err != nil {
		return nil, err
	}

	// Register default translations
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	// Register custom translations
	if err := registerCustomTranslations(validate, trans); err != nil {
		return nil, err
	}

	return &Validator{
		validate: validate,
		trans:    trans,
	}, nil
}

func (v *Validator) Validate(s interface{}) ([]ValidationError, error) {
	if err := v.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.formatErrors(validationErrors), err
		}

		return nil, err
	}

	return nil, nil
}

func (v *Validator) formatErrors(errors validator.ValidationErrors) []ValidationError {
	var validationErrors []ValidationError
	for _, err := range errors {
		validationErrors = append(validationErrors, ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Error: err.Translate(v.trans),
		})
	}
	return validationErrors
}

func registerCustomValidations(v *validator.Validate) error {
	validators := []struct {
		tag       string
		validator func(fl validator.FieldLevel) bool
	}{
		{
			tag:       "name",
			validator: validateName,
		},
		{
			tag:       "id",
			validator: validateID,
		},
		{
			tag:       "name_format",
			validator: validateNameFormat,
		},
		{
			tag:       "name_no_special",
			validator: validateNameNoSpecial,
		},
	}

	for _, val := range validators {
		if err := v.RegisterValidation(val.tag, val.validator); err != nil {
			return err
		}
	}
	return nil
}

func registerCustomTranslations(v *validator.Validate, trans ut.Translator) error {
	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "name",
			translation: "{0} must be between 2 and 100 characters",
		},
		{
			tag:         "id",
			translation: "{0} must be a positive number",
		},
		{
			tag:         "name_format",
			translation: "{0} contains invalid characters",
		},
		{
			tag:         "name_no_special",
			translation: "{0} must not contain special characters",
		},
		{
			tag:         "required",
			translation: "{0} is a required field",
		},
	}

	for _, t := range translations {
		if err := v.RegisterTranslation(t.tag, trans, func(ut ut.Translator) error {
			return ut.Add(t.tag, t.translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(fe.Tag(), fe.Field())
			return t
		}); err != nil {
			return err
		}
	}
	return nil
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	name = strings.TrimSpace(name)
	return len(name) >= minNameLength && len(name) <= maxNameLength
}

func validateID(fl validator.FieldLevel) bool {
	id := fl.Field().Int()
	return id >= minIDValue
}

func validateNameFormat(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	matched, _ := regexp.MatchString(nameRegexString, name)
	return matched
}

func validateNameNoSpecial(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	matched, _ := regexp.MatchString("^[a-zA-Z0-9\\s]+$", name)
	return matched
}
