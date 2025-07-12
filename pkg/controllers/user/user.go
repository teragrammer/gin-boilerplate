package user

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

type ManageUserController struct {
	h configs.BootHandlers
}

func NewUserController(h configs.BootHandlers) *ManageUserController {
	return &ManageUserController{h: h}
}

func (controller *ManageUserController) Browse(c *gin.Context) {
	_, pageSize, offset := utilities.Paginate(c)

	// Build query conditions based on parameters
	query := controller.h.DB

	var roleId = c.Query("role_id")
	if roleId != "" {
		query = query.Where("role_id = ?", roleId)
	}

	var status = c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// other search params or filters
	var search = c.Query("search")
	if search != "" {
		var searchTerm = strings.ToLower("%" + search + "%")
		query = query.Where("LOWER(first_name) LIKE ? OR "+
			"LOWER(middle_name) LIKE ? OR "+
			"LOWER(last_name)) LIKE ? OR "+
			"LOWER(username)) LIKE ?"+
			"LOWER(phone)) LIKE ?"+
			"LOWER(email)) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	var users []migration.User
	if err := query.
		Omit("password").
		Offset(offset).Limit(pageSize).
		Find(&users).Error; err != nil {
		utilities.LogWithLine("pkg.controllers.user.user.Browse", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (controller *ManageUserController) View(c *gin.Context) {
	// Build query conditions based on parameters
	query := controller.h.DB

	var user migration.User
	if err := query.
		Omit("password").
		Where("id = ?", c.Param("id")).
		First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    configs.Errors().E9.Code,
			"message": configs.Errors().E9.Message,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *ManageUserController) Create(c *gin.Context) {
	var form struct {
		FirstName       string  `form:"first_name" validate:"required,min=1,max=100" json:"first_name"`
		MiddleName      *string `form:"middle_name" validate:"omitempty,max=100" json:"middle_name"`
		LastName        *string `form:"last_name" validate:"required,min=1,max=100" json:"last_name"`
		Gender          *string `form:"gender" validate:"omitempty,oneof=Male Female" json:"gender"`
		Address         *string `form:"address" validate:"omitempty,max=100" json:"address"`
		BirthDate       *string `form:"birth_date" validate:"omitempty,datetime=2006-01-02" json:"birth_date"`
		Phone           *string `form:"phone" validate:"omitempty,phone=PH" json:"phone"`
		IsPhoneVerified *int    `form:"is_phone_verified" validate:"required,oneof=0 1" json:"is_phone_verified"`
		Email           *string `form:"email" validate:"omitempty,email" json:"email"`
		IsEmailVerified *int    `form:"is_email_verified" validate:"required,oneof=0 1" json:"is_email_verified"`
		RoleId          uint    `form:"role_id" validate:"required,numeric" json:"role_id"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if phone already exists
	var _phone migration.User
	if err := controller.h.DB.Where("phone", form.Phone).First(&_phone).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    configs.Errors().E5.Code,
			"message": configs.Errors().E5.Message,
		})
		return
	}

	var isPhoneVerified = *form.IsPhoneVerified == 1
	var isEmailVerified = *form.IsEmailVerified == 1

	var role = migration.User{
		FirstName:       form.FirstName,
		MiddleName:      utilities.ValueOfNullString(form.MiddleName),
		LastName:        utilities.ValueOfNullString(form.LastName),
		Gender:          utilities.ValueOfNullString(form.Gender),
		Address:         utilities.ValueOfNullString(form.Address),
		BirthDate:       utilities.ParseValueOfNullNullTime(form.BirthDate, "2006-01-02"),
		Phone:           utilities.ValueOfNullString(form.Phone),
		IsPhoneVerified: &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isPhoneVerified}},
		Email:           utilities.ValueOfNullString(form.Email),
		IsEmailVerified: &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: isEmailVerified}},
		RoleId:          form.RoleId,
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

func (controller *ManageUserController) Update(c *gin.Context) {
	var form struct {
		Name        string  `form:"name" validate:"required,min=1,max=100" json:"name"`
		Description *string `form:"description" validate:"omitempty,min=1,max=200" json:"description"`
		Slug        string  `form:"slug" validate:"required,min=1,max=100" json:"slug"`
		Rank        uint    `form:"rank" validate:"required,numeric" json:"rank"`
		IsDisabled  *int    `form:"is_disabled" validate:"required,oneof=0 1" json:"is_disabled"`
		IsActive    *int    `form:"is_active" validate:"required,oneof=0 1" json:"is_active"`
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

	var isDisabled = *form.IsDisabled == 1
	var isActive = *form.IsActive == 1

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

func (controller *ManageUserController) Delete(c *gin.Context) {
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
