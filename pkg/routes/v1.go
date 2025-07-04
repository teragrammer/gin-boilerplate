package routes

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/pkg/controllers/application"
	"gin-boilerplate/pkg/controllers/authentication"
	"gin-boilerplate/pkg/controllers/user"
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

		api.POST("/password-recovery/send", middlewares.ApplicationKeyMiddleware(h), authentication.NewPasswordRecoveryController(h).Send)
		api.POST("/password-recovery/validate", middlewares.ApplicationKeyMiddleware(h), authentication.NewPasswordRecoveryController(h).Validate)

		api.PATCH("/account/information", middlewares.ApplicationKeyMiddleware(h), middlewares.AuthenticateTokenMiddleware(h, true), user.NewAccountController(h).Information)
		api.PATCH("/account/password", middlewares.ApplicationKeyMiddleware(h), middlewares.AuthenticateTokenMiddleware(h, true), user.NewAccountController(h).Password)

		api.GET("/settings", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).Browse)
		api.GET("/settings/values", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).Values)
		api.GET("/settings/:id", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).View)
		api.POST("/settings", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).Create)
		api.PATCH("/settings", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).Update)
		api.DELETE("/settings/:id", middlewares.ApplicationKeyMiddleware(h),
			middlewares.AuthenticateTokenMiddleware(h, true), application.NewSettingController(h).Delete)
	}
}
