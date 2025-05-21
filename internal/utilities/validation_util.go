package utilities

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/nyaruka/phonenumbers"
	"regexp"
)

type ExtendValidator struct {
	validator *validator.Validate
}

// NewExtendedValidator initializes a new ExtendingValidator instance
func NewExtendedValidator() *ExtendValidator {
	v := validator.New(validator.WithRequiredStructEnabled())

	cv := &ExtendValidator{validator: v}

	// Register custom validation functions
	if err := cv.registerCustomValidations(); err != nil {
		panic(fmt.Sprintf("error registering custom validators: %s", err))
	}

	return cv
}

// registerCustomValidations registers custom validation functions
func (cv *ExtendValidator) registerCustomValidations() error {
	// check if value is a valid phone number
	if err := cv.validator.RegisterValidation("phone", cv.phoneValidation); err != nil {
		return err
	}

	// check if password is complex
	if err := cv.validator.RegisterValidation("password", cv.passwordValidation); err != nil {
		return err
	}

	// Add more custom validations as needed

	return nil
}

func (cv *ExtendValidator) GetTranslation(lang string) ut.Translator {
	// Create a new universal translator and configure it with English
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator(lang)

	// Register translation for the phone validation in English
	err := cv.validator.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} is not a valid format of phone number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
	if err != nil {
		return nil
	}

	err = cv.validator.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} is not a valid format of password", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	if err != nil {
		return nil
	}

	// Add more custom validation error messages

	_ = entranslations.RegisterDefaultTranslations(cv.validator, trans)

	return trans
}

// Validate validates the struct and returns error if validation fails
func (cv *ExtendValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// phoneValidation is a custom validation function
func (cv *ExtendValidator) phoneValidation(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	param := fl.Param()

	parsedNumber, err := phonenumbers.Parse(fieldValue, param)
	if err != nil {
		return false
	}
	isValid := phonenumbers.IsValidNumber(parsedNumber)
	if !isValid {
		return false
	}

	return true
}

func (cv *ExtendValidator) passwordValidation(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	regexLower := regexp.MustCompile("[a-z]+")
	regexDigit := regexp.MustCompile("[0-9]+")

	return regexLower.MatchString(fieldValue) &&
		regexDigit.MatchString(fieldValue)
}
