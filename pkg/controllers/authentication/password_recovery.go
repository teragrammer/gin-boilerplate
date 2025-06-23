package authentication

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PasswordRecoveryController struct {
	h configs.BootHandlers
}

func NewPasswordRecoveryController(h configs.BootHandlers) *PasswordRecoveryController {
	return &PasswordRecoveryController{h: h}
}

func (controller *PasswordRecoveryController) Send(c *gin.Context) {
	var form struct {
		To    string  `form:"to" validate:"required,oneof:email phone" json:"to"`
		Email *string `form:"email" validate:"omitempty,email" json:"email"`
		Phone *string `form:"phone" validate:"omitempty,phone=PH" json:"phone"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	var extracted = extractValueTo(form.Email, form.Phone)
	if !extracted.status {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E4.Code,
			"message": configs.Errors().E4.Message,
		})
		return
	}
}

func (controller *PasswordRecoveryController) Validate(c *gin.Context) {
	// TODO
}

type sendColumn struct {
	status bool
	name   *string
	value  *string
}

func extractValueTo(email *string, phone *string) sendColumn {
	if email != nil {
		var emailName = "email"
		return sendColumn{
			status: true,
			name:   &emailName,
			value:  email,
		}
	}

	if phone != nil {
		var phoneName = "phone"
		return sendColumn{
			status: true,
			name:   &phoneName,
			value:  phone,
		}
	}

	return sendColumn{
		status: false,
	}
}
