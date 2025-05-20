package handlers

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ErrorKeyValuePair struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationHandler(c *gin.Context, obj any) gin.H {
	if err := c.ShouldBind(obj); err != nil {
		return gin.H{
			"code":    configs.Errors().E3.Code,
			"message": err.Error(),
		}
	}

	validate := utilities.NewExtendedValidator()
	translation := validate.GetTranslation("en")
	errValidate := validate.Validate(obj)

	if errValidate != nil {
		var e []interface{}

		for _, err := range errValidate.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(obj).Elem().FieldByName(err.Field())

			e = append(e, ErrorKeyValuePair{
				Field:   field.Tag.Get("json"),
				Message: err.Translate(translation),
			})
		}

		return gin.H{
			"code":    configs.Errors().E4.Code,
			"message": configs.Errors().E4.Message,
			"errors":  e,
		}
	}

	return nil
}

func ErrorHandler(field string, message string) gin.H {
	return gin.H{
		"code":    configs.Errors().E4.Code,
		"message": configs.Errors().E4.Message,
		"errors": [...]ErrorKeyValuePair{{
			Field:   field,
			Message: message,
		}},
	}
}
