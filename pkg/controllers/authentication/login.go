package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/middlewares"
	"gin-boilerplate/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(h configs.BootHandlers) func(c *gin.Context) {
	settingSlugs := []string{
		"mx_log_try", "lck_prd",
		"tkn_lth", "tkn_exp", "tfa_req",
	}

	type Form struct {
		Username string `form:"username" validate:"required" json:"username"`
		Password string `form:"password" validate:"required" json:"password"`
	}

	return func(c *gin.Context) {
		var form Form
		e := handlers.ValidationHandler(c, &form)
		if e != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
			return
		}

		var user migration.User
		if err := h.DB.Where("username", form.Username).Preload("Role").First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, handlers.ErrorHandler("username", "The provided username is not valid, or the account does not exist"))
			return
		}

		if user.FailedLoginExpiredAt != nil {
			loginExpiredAt := user.FailedLoginExpiredAt.Time.Unix()
			currentTime := time.Now().Unix()

			if loginExpiredAt >= currentTime {
				// check if failed login tries exceed
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    configs.Errors().E10.Code,
					"message": configs.Errors().E10.Message,
				})
				return
			} else {
				// reset failed login expiration
				if err := h.DB.
					Where("id = ?", user.Id).
					Updates(migration.User{
						FailedLoginExpiredAt: &utilities.NullTime{NullTime: sql.NullTime{Time: time.Now(), Valid: true}},
					}).Error; err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code":    configs.Errors().E7.Code,
						"message": configs.Errors().E7.Message,
					})
					return
				}
			}
		}

		settings, err := repositories.Settings(h.DB, settingSlugs)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    configs.Errors().E8.Code,
				"message": configs.Errors().E8.Message,
			})
			return
		}

		_, err = utilities.VerifyHash(form.Password, user.Password)
		if err != nil {
			// increase the login attempt failed
			h.DB.Exec("UPDATE users SET login_tries = login_tries + 1 WHERE id = ?", user.Id)

			totalLoginTries := user.LoginTries + 1
			if totalLoginTries >= uint(settings.MxLogTry) {
				h.DB.Model(&migration.User{}).Where("id = ?", user.Id).
					Updates(migration.User{FailedLoginExpiredAt: &utilities.NullTime{NullTime: sql.NullTime{Time: utilities.AddMinute(time.Now(), settings.LckPrd), Valid: true}}})

				// too many login attempts
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    configs.Errors().E11.Code,
					"message": configs.Errors().E11.Message,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("password", "The username and password do not match"))
			return
		}

		// add role information
		if err := h.DB.Where("id", user.RoleId).First(&user.Role).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": configs.Errors().E9.Code, "message": configs.Errors().E9.Message})
			return
		}

		// generate authentication token
		token, err := repositories.GenerateToken(h.DB, settings, user.Id)
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
}

func Logout(h configs.BootHandlers) func(c *gin.Context) {
	return func(c *gin.Context) {
		authentication, _ := c.Get("credential")
		var id = authentication.(middlewares.Credential).Token.Id

		var result = h.DB.
			Where("id = ?", id).
			Delete(&migration.AuthenticationToken{})
		if err := result.Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    configs.Errors().E7.Code,
				"message": configs.Errors().E7.Message,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": result.RowsAffected,
		})
	}
}
