package authentication

import (
	"database/sql"
	configs "gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/handlers"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(h configs.BootHandlers) func(c *gin.Context) {
	settingSlugs := []string{"tkn_lth", "tkn_exp", "tfa_req"}

	return func(c *gin.Context) {
		type Form struct {
			FirstName  string  `form:"first_name" validate:"required,max=100" json:"first_name"`
			MiddleName *string `form:"middle_name" validate:"omitempty,max=100" json:"middle_name"`
			LastName   string  `form:"last_name" validate:"required,max=100" json:"last_name"`
			Email      *string `form:"email" validate:"omitempty,email" json:"email"`
			Phone      *string `form:"phone" validate:"omitempty,phone=PH" json:"phone"`
			Username   string  `form:"username" validate:"required,max=16,alphanum" json:"username"`
			Password   string  `form:"password" validate:"required,max=28,min=6,password" json:"password"`
		}

		var form Form
		e := handlers.ValidationHandler(c, &form)
		if e != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, e)
			return
		}

		settings, err := repositories.Settings(h.DB, settingSlugs)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    configs.Errors().E8.Code,
				"message": configs.Errors().E8.Message,
			})
			return
		}

		// check if username already exists
		var user migration.User
		if err := h.DB.Where("username", form.Username).First(&user).Error; err == nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E5.Code,
				"message": configs.Errors().E5.Message,
			})
			return
		}

		// check if email already exists
		if form.Email != nil {
			if err := h.DB.Where("email", form.Email).First(&user).Error; err == nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("email", "Email already taken"))
				return
			}
		}

		// check if phone already exists
		if form.Phone != nil {
			if err := h.DB.Where("phone", form.Phone).First(&user).Error; err == nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("email", "Phone number already taken"))
				return
			}
		}

		// default role for newly registered user
		var role migration.Role
		if err := h.DB.Where("slug", "customer").First(&role).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    configs.Errors().E9.Code,
				"message": configs.Errors().E9.Message,
			})
			return
		}

		// hash password
		hash, err := utilities.Hash(form.Password + h.Env.Security.HashSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, handlers.ErrorHandler("password", configs.Errors().E5.Message))
			return
		}

		// register the new account
		user = migration.User{
			FirstName:  form.FirstName,
			MiddleName: utilities.ValueOfNullString(form.MiddleName),
			LastName:   &utilities.NullString{NullString: sql.NullString{String: form.LastName, Valid: true}},
			RoleId:     role.Id,
			Email:      utilities.ValueOfNullString(form.Email),
			Phone:      utilities.ValueOfNullString(form.Phone),
			Username:   form.Username,
			Password:   hash,
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E7.Code,
				"message": configs.Errors().E7.Message,
			})
			return
		}

		// generate authentication token
		token, err := repositories.GenerateToken(h.DB, settings, user.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    configs.Errors().E7.Code,
				"message": configs.Errors().E7.Message,
			})
			return
		}

		// add role information
		user.Role = &role

		c.JSON(http.StatusOK, gin.H{
			"user":       user,
			"credential": token,
		})
	}
}
