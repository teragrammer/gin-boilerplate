package middlewares

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Credential struct {
	Token migration.AuthenticationToken
	User  migration.User
}

func AuthenticateTokenMiddleware(h configs.BootHandlers, halt bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedToken := c.Request.Header.Get("Authorization")
		if encodedToken == "" && !halt {
			c.Next()
			return
		}

		decodedToken, err := utilities.DecodeBase64URL(encodedToken)
		if err != nil {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E12.Code,
				"message": configs.Errors().E12.Message,
			})
			return
		}

		parts := strings.Split(string(decodedToken), ".")
		if len(parts) != 4 {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E13.Code,
				"message": configs.Errors().E13.Message,
			})
			return
		}

		var user = migration.User{}
		if err = h.DB.Where("id = ?", parts[0]).Preload("Role").First(&user).Error; err != nil {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E14.Code,
				"message": configs.Errors().E14.Message,
			})
			return
		}

		var authenticationToken = migration.AuthenticationToken{}
		if err = h.DB.Where("id = ?", parts[1]).First(&authenticationToken).Error; err != nil {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E15.Code,
				"message": configs.Errors().E15.Message,
			})
			return
		}

		if authenticationToken.Token != parts[2] {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E16.Code,
				"message": configs.Errors().E16.Message,
			})
			return
		}

		unixTimeExpiration, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E17.Code,
				"message": configs.Errors().E17.Message,
			})
			return
		}

		now := time.Now().Unix()
		if unixTimeExpiration < now ||
			(authenticationToken.ExpiredAt != nil && authenticationToken.ExpiredAt.Valid == true && authenticationToken.ExpiredAt.Time.Unix() < now) {
			if !halt {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E18.Code,
				"message": configs.Errors().E18.Message,
			})
			return
		}

		c.Set("credential", Credential{
			Token: authenticationToken,
			User:  user,
		})

		c.Next()
	}
}
