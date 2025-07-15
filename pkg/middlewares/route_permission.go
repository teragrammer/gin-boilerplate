package middlewares

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoutePermission struct {
	IsPermitted bool
}

func RoutePermissionMiddleware(roles []string, halt bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		credential, exists := c.Get("credential")
		var allow = false

		if !exists && halt {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E30.Code,
				"message": configs.Errors().E30.Message,
			})
			return
		} else if !exists && !halt {
			allow = true
		}

		if exists {
			var roleSlug = credential.(Credential).User.Role.Slug
			var isRoleExists = utilities.IsStringValueExistOnArray(roles, &roleSlug)
			if !isRoleExists && halt {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    configs.Errors().E30.Code,
					"message": configs.Errors().E30.Message,
				})
				return
			} else if !isRoleExists && !halt {
				allow = true
			}
		}

		if !allow {
			c.Set("route.permission", RoutePermission{
				IsPermitted: true,
			})
		}

		c.Next()
	}
}
