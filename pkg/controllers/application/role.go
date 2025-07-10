package application

import (
	"database/sql"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RoleController struct {
	h configs.BootHandlers
}

func NewRoleController(h configs.BootHandlers) *RoleController {
	return &RoleController{h: h}
}

func (controller *RoleController) Browse(c *gin.Context) {
	_, pageSize, offset := utilities.Paginate(c)

	// Build query conditions based on parameters
	query := controller.h.DB

	var isPublic = c.Query("is_public")
	if isPublic != "" {
		query = query.Where("is_public = ?", isPublic)
	}

	var isActive = c.Query("is_active")
	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	var slug = c.Query("slug")
	if slug != "" {
		query = query.Where("slug = ?", slug)
	}

	// other search params or filters
	var search = c.Query("search")
	if search != "" {
		var searchTerm = strings.ToLower("%" + search + "%")
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(description)) LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	var roles []migration.Role
	if err := query.
		Offset(offset).Limit(pageSize).
		Find(&roles).Error; err != nil {
		utilities.LogWithLine("pkg.controllers.application.role.Browse", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (controller *RoleController) View(c *gin.Context) {
	// Build query conditions based on parameters
	query := controller.h.DB

	var role migration.Role
	if err := query.
		Where("id = ?", c.Param("id")).
		First(&role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    configs.Errors().E9.Code,
			"message": configs.Errors().E9.Message,
		})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (controller *RoleController) Create(c *gin.Context) {
	var form struct {
		Name        string  `form:"name" validate:"required,min=1,max=100" json:"name"`
		Description *string `form:"description" validate:"omitempty,min=1,max=200" json:"description"`
		Slug        string  `form:"slug" validate:"required,min=1,max=100" json:"slug"`
		Rank        uint    `form:"rank" validate:"required,numeric,regexp=^[0-9]+$" json:"rank"`
		IsDisabled  int     `form:"is_disabled" validate:"required,oneof=0 1" json:"is_disabled"`
		IsActive    int     `form:"is_active" validate:"required,oneof=0 1" json:"is_active"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if slug already exists
	var _role migration.Role
	if err := controller.h.DB.Where("slug", form.Slug).First(&_role).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E5.Code,
			"message": configs.Errors().E5.Message,
		})
		return
	}

	var isDisabled = form.IsDisabled == 1
	var isActive = form.IsActive == 1

	var role = migration.Role{
		Name:        form.Name,
		Description: utilities.ValueOfNullString(form.Description),
		Slug:        form.Slug,
		Rank:        form.Rank,
		IsPublic:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isDisabled}},
		IsActive:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isActive}},
	}

	if err := controller.h.DB.
		Create(&role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": role.Id,
	})
}

func (controller *RoleController) Update(c *gin.Context) {
	var form struct {
		Name        string  `form:"name" validate:"required,min=1,max=100" json:"name"`
		Description *string `form:"description" validate:"omitempty,min=1,max=200" json:"description"`
		Slug        string  `form:"slug" validate:"required,min=1,max=100" json:"slug"`
		Rank        uint    `form:"rank" validate:"required,numeric,regexp=^[0-9]+$" json:"rank"`
		IsDisabled  int     `form:"is_disabled" validate:"required,oneof=0 1" json:"is_disabled"`
		IsActive    int     `form:"is_active" validate:"required,oneof=0 1" json:"is_active"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if slug already exists
	id := c.Param("id")
	var _role migration.Role
	if err := controller.h.DB.
		Where("id <> ?", id).
		Where("slug", form.Slug).
		First(&_role).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E5.Code,
			"message": configs.Errors().E5.Message,
		})
		return
	}

	var isDisabled = form.IsDisabled == 1
	var isActive = form.IsActive == 1

	var role = migration.Role{
		Name:        form.Name,
		Description: utilities.ValueOfNullString(form.Description),
		Slug:        form.Slug,
		Rank:        form.Rank,
		IsPublic:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isDisabled}},
		IsActive:    &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isActive}},
	}
	if err := controller.h.DB.
		Where("id = ?", id).
		Updates(role).Error; err != nil {
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

func (controller *RoleController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := controller.h.DB.
		Where("id = ?", id).
		Delete(&migration.Role{}).Error; err != nil {
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
