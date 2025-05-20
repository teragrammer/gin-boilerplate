package database

import (
	"fmt"
	"gin-boilerplate/configs"
	Seeding "gin-boilerplate/database/seed"
	"gorm.io/gorm"
)

func RunSeeder() {
	db := Connector(configs.Config("env.json").Db)

	seeders := []func(*gorm.DB) error{
		Seeding.SeedingSettingTable,
		Seeding.SeedingRoleTable,
		Seeding.SeedingUserTable,
	}

	// Execute each migration step sequentially
	for _, seed := range seeders {
		err := seed(db)
		if err != nil {
			return
		}
	}

	fmt.Println("Seeding Tables Done")
}
