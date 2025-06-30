package authentication

import (
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

}
