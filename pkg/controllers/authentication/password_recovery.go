package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

const MaxTries = 5
const NextTryMinutes = 3

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
	var form struct {
		To    string  `form:"to" validate:"required,oneof=email phone" json:"to"`
		Code  string  `form:"code" validate:"max=6,min=6" json:"code"`
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

	if recovery.Tries >= MaxTries {
		var isExceededTries = true
		if recovery.NextTryAt != nil &&
			recovery.NextTryAt.Valid &&
			recovery.NextTryAt.Time.Unix() <= time.Now().Unix() {
			isExceededTries = false
		} else {
			controller.h.DB.Model(&migration.PasswordRecovery{}).Where("id = ?", recovery.Id).
				Updates(migration.PasswordRecovery{
					NextTryAt: &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: utilities.AddMinute(time.Now(), NextTryMinutes)}},
				})
		}

		if isExceededTries {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E27.Code,
				"message": configs.Errors().E27.Message,
			})
			return
		}
	}

	verified, err := utilities.VerifyHash(form.Code, recovery.Code)
	if err != nil || !verified {
		// record number of tries
		controller.h.DB.Model(&migration.PasswordRecovery{}).
			Where("id = ?", recovery.Id).
			UpdateColumn("tries", gorm.Expr("tries + 1"))

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E28.Code,
			"message": configs.Errors().E28.Message,
		})
		return
	}

	// remove all recovery
	controller.h.DB.Where("send_to = ?", *extracted.value).Delete(&migration.PasswordRecovery{})

	// set user authentication token
	var user migration.User
	if err := controller.h.DB.Where(*extracted.name, *extracted.value).First(&user).Error; err != nil {
		var fieldName = *extracted.name
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler(fieldName, "Unable to find "+fieldName))
		return
	}

	// add role information
	if err := controller.h.DB.Where("id", user.RoleId).First(&user.Role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": configs.Errors().E9.Code, "message": configs.Errors().E9.Message})
		return
	}

	// get the settings for generating authentication token
	settings, err := repositories.Settings(controller.h.DB, []string{
		"tkn_lth", "tkn_exp", "tfa_req",
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E8.Code,
			"message": configs.Errors().E8.Message,
		})
		return
	}

	// generate authentication token
	token, err := repositories.GenerateToken(controller.h.DB, settings, user.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":       user,
		"credential": token,
	})
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
