package testutilities

import (
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg/repositories"
)

func GenerateAuthentication(h configs.BootHandlers, slug string, username string) (migration.User, *migration.AuthenticationToken) {
	h.DB.Where("username = ?", username).Delete(&migration.User{})

	var role migration.Role
	h.DB.Where("slug", slug).First(&role)

	hash, _ := utilities.Hash("123456" + h.Env.Security.HashSecret)

	email, _ := utilities.GenerateRandomString(16)
	finalEmail := email + "@gmail.com"
	user := migration.User{
		RoleId:   role.Id,
		Email:    utilities.ValueOfNullString(&finalEmail),
		Username: username,
		Password: hash,
	}
	h.DB.Create(&user)

	settings, _ := repositories.Settings(h.DB, []string{
		"tkn_lth", "tkn_exp", "tfa_req",
	})

	token, _ := repositories.GenerateToken(h.DB, settings, user.Id)

	return user, token
}
