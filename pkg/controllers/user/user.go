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
		Username        string  `form:"username" validate:"required,max=32,alphanum" json:"username"`
		Password        string  `form:"password" validate:"required,min=6,max=28,password" json:"password"`
		RoleId          uint    `form:"role_id" validate:"required,numeric" json:"role_id"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	// check if phone already exists
	var _phone migration.User
	if form.Phone != nil {
		if err := controller.h.DB.Where("phone", *form.Phone).First(&_phone).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("phone", configs.Errors().E5.Message))
			return
		}
	}

	// check if email already exists
	var _email migration.User
	if form.Email != nil {
		if err := controller.h.DB.Where("email", *form.Email).First(&_email).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("email", configs.Errors().E5.Message))
			return
		}
	}

	// check if username already exists
	var _username migration.User
	if err := controller.h.DB.Where("username", form.Username).First(&_username).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("username", configs.Errors().E5.Message))
		return
	}

	// check if role exists
	var _role migration.Role
	if err := controller.h.DB.Where("id", form.RoleId).First(&_role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, handlers.ErrorHandler("role_id", configs.Errors().E9.Message))
		return
	}

	// hash password
	hash, err := utilities.Hash(form.Password + controller.h.Env.Security.HashSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("password", configs.Errors().E5.Message))
		return
	}

	var isPhoneVerified = *form.IsPhoneVerified == 1
	var isEmailVerified = *form.IsEmailVerified == 1

	var user = migration.User{
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
		Username:        form.Username,
		Password:        hash,
		RoleId:          form.RoleId,
	}

	if err := controller.h.DB.
		Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    configs.Errors().E7.Code,
			"message": configs.Errors().E7.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.Id,
	})
}

func (controller *ManageUserController) Update(c *gin.Context) {
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
		Username        string  `form:"username" validate:"required,max=32,alphanum" json:"username"`
		Password        *string `form:"password" validate:"omitempty,min=6,max=28,password" json:"password"`
		RoleId          uint    `form:"role_id" validate:"required,numeric" json:"role_id"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	id := c.Param("id")

	// check if phone already exists
	var _phone migration.User
	if form.Phone != nil {
		if err := controller.h.DB.Where("phone", *form.Phone).
			Where("id <> ?", id).
			First(&_phone).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("phone", configs.Errors().E5.Message))
			return
		}
	}

	// check if email already exists
	var _email migration.User
	if form.Email != nil {
		if err := controller.h.DB.Where("email", *form.Email).
			Where("id <> ?", id).
			First(&_email).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("email", configs.Errors().E5.Message))
			return
		}
	}

	// check if username already exists
	var _username migration.User
	if err := controller.h.DB.Where("username", form.Username).
		Where("id <> ?", id).
		First(&_username).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("username", configs.Errors().E5.Message))
		return
	}

	// check if role exists
	var _role migration.Role
	if err := controller.h.DB.Where("id", form.RoleId).
		First(&_role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, handlers.ErrorHandler("role_id", configs.Errors().E9.Message))
		return
	}

	var isPhoneVerified = *form.IsPhoneVerified == 1
	var isEmailVerified = *form.IsEmailVerified == 1

	var user = migration.User{
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
		Username:        form.Username,
		RoleId:          form.RoleId,
	}

	// hash password
	if form.Password != nil {
		hash, err := utilities.Hash(*form.Password + controller.h.Env.Security.HashSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("password", configs.Errors().E5.Message))
			return
		}

		user.Password = hash
	}

	if err := controller.h.DB.
		Where("id = ?", id).
		Updates(user).Error; err != nil {
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
		Delete(&migration.User{}).Error; err != nil {
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
