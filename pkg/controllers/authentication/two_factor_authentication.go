package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TFAController struct {
	h configs.BootHandlers
}

func NewTFAController(h configs.BootHandlers) *TFAController {
	return &TFAController{h: h}
}

func (controller *TFAController) Send(c *gin.Context) {
	credential, _ := c.Get("credential")

	// check if tfa is required
	// to save resources
	if credential.(middlewares.Credential).Token.IsTFARequired != nil && credential.(middlewares.Credential).Token.IsTFARequired.Valid &&
		!credential.(middlewares.Credential).Token.IsTFARequired.Bool {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E19.Code,
			"message": configs.Errors().E19.Message,
		})
		return
	}

	// check if user has valid email
	if credential.(middlewares.Credential).User.Email == nil ||
		credential.(middlewares.Credential).User.Email.Valid == false {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E20.Code,
			"message": configs.Errors().E20.Message,
		})
		return
	}

	// generate the random code for tfa
	code, err := utilities.GenerateRandomNumber(8)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	var nextTryAt = utilities.AddMinute(time.Now(), 2)
	var tfa migration.TwoFactorAuthentication
	if err := controller.h.DB.Where("token_id", credential.(middlewares.Credential).Token.Id).First(&tfa).Error; err != nil {
		tfa = migration.TwoFactorAuthentication{
			TokenId:    credential.(middlewares.Credential).Token.Id,
			Code:       code,
			ExpiredAt:  &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: utilities.AddMinute(time.Now(), 5)}},
			NextSendAt: &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: nextTryAt}},
		}
		if err := controller.h.DB.Create(&tfa).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E7.Code,
				"message": configs.Errors().E7.Message,
			})
			return
		}
	} else if tfa.NextSendAt != nil && tfa.NextSendAt.Valid {
		if time.Now().Unix() < tfa.NextSendAt.Time.Unix() {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E21.Code,
				"message": configs.Errors().E21.Message,
			})
			return
		}

		controller.h.DB.Model(&migration.TwoFactorAuthentication{}).Where("id = ?", tfa.Id).
			Updates(migration.TwoFactorAuthentication{
				Code:       code,
				ExpiredAt:  &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: utilities.AddMinute(time.Now(), 5)}},
				NextSendAt: &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: nextTryAt}},
			})
	}

	if controller.h.Env.Environment == "production" {
		// TODO
		// send email with tfa code
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       tfa.Id,
		"next_try": nextTryAt,
	})
}

func (controller *TFAController) Validate(c *gin.Context) {
	var form struct {
		Code string `form:"code" validate:"required" json:"code"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	credential, _ := c.Get("credential")

	// check if tfa is required
	// to save resources
	if credential.(middlewares.Credential).Token.IsTFARequired != nil && credential.(middlewares.Credential).Token.IsTFARequired.Valid &&
		!credential.(middlewares.Credential).Token.IsTFARequired.Bool {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E19.Code,
			"message": configs.Errors().E19.Message,
		})
		return
	}

	var tfa migration.TwoFactorAuthentication
	if err := controller.h.DB.Where("token_id", credential.(middlewares.Credential).Token.Id).First(&tfa).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E9.Code,
			"message": configs.Errors().E9.Message,
		})
		return
	}

	// check for expiration
	if tfa.ExpiredAt != nil && tfa.ExpiredAt.Valid {
		if time.Now().Unix() > tfa.ExpiredAt.Time.Unix() {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E22.Code,
				"message": configs.Errors().E22.Message,
			})
			return
		}
	}

	if tfa.ExpiredTriesAt != nil && tfa.ExpiredTriesAt.Valid {
		// multiple pending tries
		if time.Now().Unix() < tfa.ExpiredTriesAt.Time.Unix() {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E23.Code,
				"message": configs.Errors().E23.Message,
			})
			return
		}

		// reset the tries
		controller.h.DB.Model(&migration.TwoFactorAuthentication{}).Where("id = ?", tfa.Id).
			Updates(migration.TwoFactorAuthentication{
				Tries:          0,
				ExpiredTriesAt: &utilities.NullTime{NullTime: sql.NullTime{Valid: false}},
			})
	}

	// too many failed tries
	if tfa.Tries > 5 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E23.Code,
			"message": configs.Errors().E23.Message,
		})
		return
	}

	verified, err := utilities.VerifyHash(form.Code, tfa.Code)
	if err != nil || !verified {
		// record number of tries
		controller.h.DB.Model(&migration.TwoFactorAuthentication{}).
			Where("id = ?", tfa.Id).
			UpdateColumn("tries", gorm.Expr("tries + 1"))

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E24.Code,
			"message": configs.Errors().E24.Message,
		})
		return
	}

	// update the authentication
	controller.h.DB.Model(&migration.AuthenticationToken{}).Where("id = ?", tfa.Id).
		Updates(migration.AuthenticationToken{
			IsTFAVerified: &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: true}},
		})

	// delete the tfa
	controller.h.DB.Where("id = ?", tfa.Id).Delete(&migration.TwoFactorAuthentication{})

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
