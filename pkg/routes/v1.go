package routes

import (
	"gin-boilerplate/configs"
	authentication2 "gin-boilerplate/pkg/controllers/authentication"
	"gin-boilerplate/pkg/middlewares"
)

var isV1RoutesInitialized = false

func V1Routes(h configs.BootHandlers) {
	if isV1RoutesInitialized == true {
		return
	}
	isV1RoutesInitialized = true

	api := h.Engine.Group("/v1")
	{
		api.POST("/register", middlewares.ApplicationKeyMiddleware(h), authentication2.Register(h))
		api.POST("/login", middlewares.ApplicationKeyMiddleware(h), authentication2.Login(h))
		api.GET("/logout", middlewares.ApplicationKeyMiddleware(h), authentication2.Logout(h))
	}
}
