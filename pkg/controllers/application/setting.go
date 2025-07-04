package application

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type SettingController struct {
	h configs.BootHandlers
}

func NewSettingController(h configs.BootHandlers) *SettingController {
	return &SettingController{h: h}
}

func (controller *SettingController) Browse(c *gin.Context) {
	_, pageSize, offset := utilities.Paginate(c)

	// Build query conditions based on parameters
	query := controller.h.DB

	var isDisabled = c.Query("is_disabled")
	if isDisabled != "" {
		query = query.Where("is_disabled = ?", isDisabled)
	}

	var isPublic = c.Query("is_public")
	if isPublic != "" {
		query = query.Where("is_public = ?", isPublic)
	}

	var _type = c.Query("type")
	if _type != "" {
		query = query.Where("type = ?", _type)
	}

	// other search params or filters
	var search = c.Query("search")
	if search != "" {
		var searchTerm = strings.ToLower("%" + search + "%")
		query = query.Where("LOWER(slug) LIKE ? OR LOWER(value) LIKE ? OR LOWER(description)) LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	var settings []migration.Setting
	if err := query.
		Offset(offset).Limit(pageSize).
		Find(&settings).Error; err != nil {
		utilities.LogWithLine("pkg.controllers.application.setting.Browse", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (controller *SettingController) Values(c *gin.Context) {
	visibility := "public"
	settings, err := repositories.Settings(controller.h.DB, []string{}, &visibility)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (controller *SettingController) View(c *gin.Context) {
	// Build query conditions based on parameters
	query := controller.h.DB

	var setting migration.Setting
	if err := query.
		Where("id = ?", c.Param("id")).
		First(&setting).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    configs.Errors().E9.Code,
			"message": configs.Errors().E9.Message,
		})
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (controller *SettingController) Create(c *gin.Context) {

}

func (controller *SettingController) Update(c *gin.Context) {

}

func (controller *SettingController) Delete(c *gin.Context) {

}
