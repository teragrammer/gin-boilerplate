package database

import (
	"fmt"
	"gin-boilerplate/configs"
	"gin-boilerplate/database/migration"
)

func RunMigrations() {
	db := Connector(configs.Config("env.json").Db)

	// Auto create/update table based on model definition
	err := db.AutoMigrate(
		&migration.Setting{},
		&migration.Role{},
		&migration.User{},
		&migration.AuthenticationToken{},
		&migration.TwoFactorAuthentication{},
	)
	if err != nil {
		fmt.Println("Failed Migrating", err.Error())
		return
	}

	fmt.Println("Migrating Tables Done")
}
