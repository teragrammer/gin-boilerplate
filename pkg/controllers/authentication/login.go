package authentication

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(h configs.BootHandlers) func(c *gin.Context) {
	settingSlugs := []string{"mx_log_try", "tkn_exp"}

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
					"code":    configs.Errors().E10,
					"message": configs.Errors().E10,
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
						"code":    configs.Errors().E7,
						"message": configs.Errors().E7,
					})
					return
				}
			}
		}

		_, err := utilities.VerifyHash(form.Password, user.Password)
		if err != nil {
			// TODO
			_, err = repositories.Settings(h.DB, settingSlugs)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    configs.Errors().E8,
					"message": configs.Errors().E8,
				})
				return
			}

			// TODO
			// increase the login attempt failed

			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("password", "The username and password do not match"))
			return
		}
	}
}

func Logout(h configs.BootHandlers) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
