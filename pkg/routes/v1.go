package routes

import (
	"gin-boilerplate/configs"
	authentication2 "gin-boilerplate/pkg/controllers/authentication"
	"gin-boilerplate/pkg/middlewares"
)

func V1Routes(h configs.BootHandlers) {
	api := h.Engine.Group("/v1")
	{
		api.POST("/register", middlewares.ApplicationKeyMiddleware(h), authentication2.Register(h))
		api.POST("/login", middlewares.ApplicationKeyMiddleware(h), authentication2.Login(h))
		api.GET("/logout", middlewares.ApplicationKeyMiddleware(h), middlewares.AuthenticateTokenMiddleware(h, true), authentication2.Logout(h))
	}
}
