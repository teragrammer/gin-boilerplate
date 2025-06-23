package routes

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/pkg/controllers/authentication"
	"gin-boilerplate/pkg/middlewares"
)

func V1Routes(h configs.BootHandlers) {
	api := h.Engine.Group("/v1")
	{
		api.POST("/register", middlewares.ApplicationKeyMiddleware(h), authentication.Register(h))
		api.POST("/login", middlewares.ApplicationKeyMiddleware(h), authentication.Login(h))
		api.GET("/logout", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), authentication.Logout(h))

		api.GET("/tfa/send", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), authentication.NewTFAController(h).Send)
		api.POST("/tfa/validate", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), authentication.NewTFAController(h).Validate)

		api.GET("/password-recovery/send", middlewares.ApplicationKeyMiddleware(h), authentication.NewPasswordRecoveryController(h).Send)
		api.POST("/password-recovery/validate", middlewares.ApplicationKeyMiddleware(h), authentication.NewPasswordRecoveryController(h).Validate)
	}
}
