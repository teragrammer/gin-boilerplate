package middlewares

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApplicationKeyMiddleware(h configs.BootHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("X-Secret-Key")

		if h.Env.App.Key != "" && key != h.Env.App.Key {
			utilities.LogWithLine("middlewares.application_key_middleware", utilities.GetClientIP(c))

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    configs.Errors().E1.Code,
				"message": configs.Errors().E1.Message,
			})
			return
		}

		c.Next()
	}
}
