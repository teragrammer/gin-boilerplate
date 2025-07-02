package authentication

import (
	"database/sql"
	"fmt"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountController struct {
	h configs.BootHandlers
}

func NewAccountController(h configs.BootHandlers) *AccountController {
	return &AccountController{h: h}
}

func (controller *AccountController) Information(c *gin.Context) {
	var form struct {
		FirstName  string  `form:"first_name" validate:"required,min=1,max=100" json:"first_name"`
		MiddleName *string `form:"middle_name" validate:"omitempty,min=1,max=100" json:"middle_name"`
		LastName   *string `form:"last_name" validate:"required,min=1,max=100" json:"last_name"`
		Address    *string `form:"address" validate:"omitempty,min=1,max=100" json:"address"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	credential, _ := c.Get("credential")
	if err := controller.h.DB.
		Where("id = ?", credential.(middlewares.Credential).User.Id).
		Updates(migration.User{
			FirstName:  form.FirstName,
			MiddleName: utilities.ValueOfNullString(form.MiddleName),
			LastName:   utilities.ValueOfNullString(form.LastName),
			Address:    utilities.ValueOfNullString(form.Address),
		}).Error; err != nil {
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

func (controller *AccountController) Password(c *gin.Context) {
	var form struct {
		CurrentPassword string  `form:"current_password" validate:"required,min=1,max=100" json:"current_password"`
		NewPassword     *string `form:"new_password" validate:"omitempty,min=6,max=32" json:"new_password"`
		Username        *string `form:"username" validate:"omitempty,min=2,max=16" json:"username"`
		Email           *string `form:"email" validate:"omitempty,min=1,max=100,email" json:"email"`
		Phone           *string `form:"phone" validate:"omitempty,min=10,max=16,phone=PH" json:"phone"`
	}

	e := handlers.ValidationHandler(c, &form)
	if e != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
		return
	}

	credential, _ := c.Get("credential")
	_, err := utilities.VerifyHash(form.CurrentPassword, credential.(middlewares.Credential).User.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    configs.Errors().E29.Code,
			"message": configs.Errors().E29.Message,
		})
		return
	}

	var _data migration.User

	// hashed if new password
	if form.NewPassword != nil {
		hash, err := utilities.Hash(*form.NewPassword + controller.h.Env.Security.HashSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("new_password", configs.Errors().E5.Message))
			return
		}

		_data.Password = hash
	}

	// check if username has no duplicate
	if form.Username != nil {
		var username migration.User
		if err := controller.h.DB.
			Where("id <> ?", credential.(middlewares.Credential).User.Id).
			Where("username = ?", *form.Username).
			First(&username).Error; err == nil {
			fmt.Print(username.Username, "=", username.Id, "=", credential.(middlewares.Credential).User.Id)
			c.AbortWithStatusJSON(http.StatusConflict, handlers.ErrorHandler("username", configs.Errors().E5.Message))
			return
		}

		_data.Username = *form.Username
	}

	// check if email has no duplicate
	if form.Email != nil {
		var email migration.User
		if err := controller.h.DB.
			Where("id <> ?", credential.(middlewares.Credential).User.Id).
			Where("email = ?", *form.Email).
			First(&email).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusConflict, handlers.ErrorHandler("email", configs.Errors().E5.Message))
			return
		}

		_data.Email = &utilities.NullString{NullString: sql.NullString{String: *form.Email, Valid: true}}
	}

	// check if phone has no duplicate
	if form.Phone != nil {
		var phone migration.User
		if err := controller.h.DB.
			Where("id <> ?", credential.(middlewares.Credential).User.Id).
			Where("phone = ?", *form.Phone).
			First(&phone).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusConflict, handlers.ErrorHandler("phone", configs.Errors().E5.Message))
			return
		}

		_data.Phone = &utilities.NullString{NullString: sql.NullString{String: *form.Phone, Valid: true}}
	}

	if err := controller.h.DB.
		Where("id = ?", credential.(middlewares.Credential).User.Id).
		Updates(_data).Error; err != nil {
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
