package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/middlewares"
	"github.com/gin-gonic/gin"
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
		credential.(middlewares.Credential).Token.IsTFARequired.Bool {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E19.Code,
			"message": configs.Errors().E19.Message,
		})
		return
	}

	// check if user has valid email
	if credential.(middlewares.Credential).User.Email != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E20.Code,
			"message": configs.Errors().E20.Message,
		})
		return
	}

	var nextTryAt = utilities.AddMinute(time.Now(), 2)
	var tfa migration.TwoFactorAuthentication
	if err := controller.h.DB.Where("token_id", credential.(middlewares.Credential).Token.Id).First(&tfa).Error; err != nil {
		// generate the random code for tfa
		code, err := utilities.GenerateRandomNumber(8)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E7.Code,
				"message": configs.Errors().E7.Message,
			})
			return
		}

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

	}
}

func (controller *TFAController) Validate(c *gin.Context) {

}
