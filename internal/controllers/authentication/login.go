package authentication

import (
	"gin-boilerplate/internal/configs"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `form:"username" validate:"required" json:"username"`
	Password string `form:"password" validate:"required" json:"password"`
}

func Login(h configs.BootHandlers) func(c *gin.Context) {
	return func(c *gin.Context) {
		
	}
}

func Logout(h configs.BootHandlers) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
