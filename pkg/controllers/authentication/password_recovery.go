package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PasswordRecoveryController struct {
	h configs.BootHandlers
}

func NewPasswordRecoveryController(h configs.BootHandlers) *PasswordRecoveryController {
	return &PasswordRecoveryController{h: h}
}

const CodeLength = 6

const NextResendMinutes = 2
const CodeExpirationMinutes = 30

func (controller *PasswordRecoveryController) Send(c *gin.Context) {
	var form struct {
		To    string  `form:"to" validate:"required,oneof=email phone" json:"to"`
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

	var user migration.User
	if err := controller.h.DB.Where(*extracted.name, *extracted.value).First(&user).Error; err != nil {
		var fieldName = *extracted.name
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler(fieldName, "Unable to find "+fieldName))
		return
	}

	var recovery migration.PasswordRecovery
	if err := controller.h.DB.Where("send_to", *extracted.value).First(&recovery).Error; err == nil {
		if recovery.NextResendAt.Unix() > time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E25.Code,
				"message": configs.Errors().E25.Message,
			})
			return
		}
	}

	var code, err = utilities.GenerateRandomString(CodeLength)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E26.Code,
			"message": configs.Errors().E26.Message,
		})
		return
	}

	hashedCode, err := utilities.Hash(code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E26.Code,
			"message": configs.Errors().E26.Message,
		})
		return
	}

	var _type = migration.Email
	if *extracted.name == "phone" {
		_type = migration.Phone
	}
	controller.h.DB.Where("send_to = ?", *extracted.value).Delete(&migration.PasswordRecovery{})
	var recoveryData = migration.PasswordRecovery{
		Type:         _type,
		SendTo:       &utilities.NullString{NullString: sql.NullString{Valid: true, String: *extracted.value}},
		Code:         hashedCode,
		NextResendAt: utilities.AddDay(time.Now(), NextResendMinutes),
		ExpiredAt:    utilities.AddDay(time.Now(), CodeExpirationMinutes),
	}
	if err := controller.h.DB.Create(&recoveryData).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E26.Code,
			"message": configs.Errors().E26.Message,
		})
		return
	}

	// TODO
	// send code to email or phone

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
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
