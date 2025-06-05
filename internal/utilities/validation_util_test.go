package utilities

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestValidatePhoneFail(t *testing.T) {
	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	var form struct {
		Phone *string `form:"phone" validate:"omitempty,phone=US" json:"phone"`
	}

	var phone = "+1"
	form.Phone = &phone

	validate := NewExtendedValidator()
	errValidate := validate.Validate(form)

	assert.Equal(t, errValidate != nil, true)
}

func TestValidatePhoneSuccess(t *testing.T) {
	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	var form struct {
		Phone *string `form:"phone" validate:"omitempty,phone=US" json:"phone"`
	}

	var phone = "+1 (310) 555-5678"
	form.Phone = &phone

	validate := NewExtendedValidator()
	errValidate := validate.Validate(form)

	assert.Equal(t, errValidate == nil, true)
}

func TestValidatePasswordFail(t *testing.T) {
	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	var form struct {
		Phone *string `form:"password" validate:"password" json:"phone"`
	}

	var phone = "123456"
	form.Phone = &phone

	validate := NewExtendedValidator()
	errValidate := validate.Validate(form)

	assert.Equal(t, errValidate != nil, true)
}

func TestValidatePasswordSuccess(t *testing.T) {
	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	var form struct {
		Phone *string `form:"password" validate:"password" json:"phone"`
	}

	var phone = "@BC123eFG"
	form.Phone = &phone

	validate := NewExtendedValidator()
	errValidate := validate.Validate(form)

	assert.Equal(t, errValidate == nil, true)
}
