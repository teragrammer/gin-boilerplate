package seed

import (
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
)

func SeedingUserTable(db *gorm.DB) error {
	hash, _ := utilities.Hash("123456")

	var admin migration.Role
	db.Where("slug", "admin").First(&admin)
	users := migration.User{FirstName: "Admin", RoleId: admin.Id, Username: "admin", Password: hash}

	// Insert records into the database
	if err := db.Create(&users).Error; err != nil {
		return err
	} else {
		return nil
	}
}
