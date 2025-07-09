package application

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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
	var form struct {
		Name        string  `form:"name" validate:"required,min=1,max=100" json:"name"`
		Slug        string  `form:"slug" validate:"required,min=1,max=100" json:"slug"`
		Value       *string `form:"value" validate:"omitempty" json:"value"`
		Description *string `form:"description" validate:"omitempty,min=1,max=200" json:"description"`
		Type        string  `form:"type" validate:"required,oneof=string integer float boolean array" json:"type"`
		IsDisabled  int     `form:"is_disabled" validate:"required,oneof=0 1" json:"is_disabled"`
		IsPublic    int     `form:"is_public" validate:"required,oneof=0 1" json:"is_public"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if slug already exists
	var _setting migration.Setting
	if err := controller.h.DB.Where("slug", form.Slug).First(&_setting).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E5.Code,
			"message": configs.Errors().E5.Message,
		})
		return
	}

	var isDisabled = form.IsPublic == 1
	var isPublic = form.IsPublic == 1

	var setting = migration.Setting{
		Name:        form.Name,
		Slug:        form.Slug,
		Value:       utilities.ValueOfNullString(form.Value),
		Description: utilities.ValueOfNullString(form.Description),
		Type:        form.Type,
		IsDisabled:  &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isDisabled}},
		IsPublic:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isPublic}},
	}

	if err := controller.h.DB.
		Create(&setting).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": setting.Id,
	})
}

func (controller *SettingController) Update(c *gin.Context) {
	var form struct {
		Name        string  `form:"name" validate:"required,min=1,max=100" json:"name"`
		Slug        string  `form:"slug" validate:"required,min=1,max=100" json:"slug"`
		Value       *string `form:"value" validate:"omitempty" json:"value"`
		Description *string `form:"description" validate:"omitempty,min=1,max=200" json:"description"`
		Type        string  `form:"type" validate:"required,oneof=string integer float boolean array" json:"type"`
		IsDisabled  int     `form:"is_disabled" validate:"required,oneof=0 1" json:"is_disabled"`
		IsPublic    int     `form:"is_public" validate:"required,oneof=0 1" json:"is_public"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if slug already exists
	id := c.Param("id")
	var _setting migration.Setting
	if err := controller.h.DB.
		Where("id <> ?", id).
		Where("slug", form.Slug).
		First(&_setting).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E5.Code,
			"message": configs.Errors().E5.Message,
		})
		return
	}

	var isDisabled = form.IsPublic == 1
	var isPublic = form.IsPublic == 1

	var setting = migration.Setting{
		Name:        form.Name,
		Slug:        form.Slug,
		Value:       utilities.ValueOfNullString(form.Value),
		Description: utilities.ValueOfNullString(form.Description),
		Type:        form.Type,
		IsDisabled:  &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isDisabled}},
		IsPublic:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isPublic}},
	}
	if err := controller.h.DB.
		Where("id = ?", id).
		Updates(setting).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func (controller *SettingController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := controller.h.DB.
		Where("id = ?", id).
		Delete(&migration.Setting{}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
