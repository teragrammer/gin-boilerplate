package application

import (
	"gin-boilerplate/configs"
	"github.com/gin-gonic/gin"
)

type SettingController struct {
	h configs.BootHandlers
}

func NewSettingController(h configs.BootHandlers) *SettingController {
	return &SettingController{h: h}
}

func (controller *SettingController) Browse(c *gin.Context) {

}

func (controller *SettingController) Values(c *gin.Context) {

}

func (controller *SettingController) View(c *gin.Context) {

}

func (controller *SettingController) Create(c *gin.Context) {

}

func (controller *SettingController) Update(c *gin.Context) {

}

func (controller *SettingController) Delete(c *gin.Context) {

}
