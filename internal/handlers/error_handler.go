package handlers

import (
	"gin-boilerplate/internal/configs"
	"github.com/gin-gonic/gin"
)

type ErrorKeyValuePair struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormValidationHandler(c *gin.Context, obj any) gin.H {
	if err := c.ShouldBind(obj); err != nil {
		return gin.H{
			"code":    configs.Errors().E3.Code,
			"message": err.Error(),
		}
	}

	// TODO

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
