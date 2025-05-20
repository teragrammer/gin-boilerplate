package seed

import (
	"database/sql"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
)

func SeedingRoleTable(db *gorm.DB) error {
	roles := []migration.Role{
		{Name: "Admin", Slug: "admin", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "Customer", Slug: "customer", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
	}

	// Insert records into the database
	if err := db.Create(&roles).Error; err != nil {
		return err
	} else {
		return nil
	}
}
