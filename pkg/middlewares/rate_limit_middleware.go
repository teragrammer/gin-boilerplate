package middlewares

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/internal/utilities"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"time"
)

func keyFunc(c *gin.Context) string {
	return utilities.GetClientIP(c)
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	utilities.LogWithLine("middlewares.rate_limit_middleware", utilities.GetClientIP(c))

	c.JSON(429, gin.H{
		"code":    configs.Errors().E2.Code,
		"message": configs.Errors().E2.Message + ", try again in " + time.Until(info.ResetTime).String(),
	})
}

func RateLimitMiddleware(boot configs.BootHandlers) gin.HandlerFunc {
	limiterStorage := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: boot.Redis,
		Rate:        time.Second * boot.Env.Rate.Window, // time in seconds (window)
		Limit:       boot.Env.Rate.Limit,                // requester per rate
	})

	return ratelimit.RateLimiter(limiterStorage, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
