package routes

import (
	"gin-boilerplate/internal/configs"
	"gin-boilerplate/internal/controllers/authentication"
	"gin-boilerplate/internal/middlewares"
)

var isV1RoutesInitialized = false

func V1Routes(h configs.BootHandlers) {
	if isV1RoutesInitialized == true {
		return
	}
	isV1RoutesInitialized = true

	api := h.Engine.Group("/v1")
	{
		api.POST("/register", middlewares.ApplicationKeyMiddleware(h), authentication.Register(h))
		api.POST("/login", middlewares.ApplicationKeyMiddleware(h), authentication.Login(h))
		api.GET("/logout", middlewares.ApplicationKeyMiddleware(h), authentication.Logout(h))
	}
}
