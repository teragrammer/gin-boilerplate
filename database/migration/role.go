package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	Id          uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string                `gorm:"type:varchar(100);not null" json:"name"`
	Description *utilities.NullString `gorm:"type:TINYTEXT" json:"description"`
	Slug        string                `gorm:"type:varchar(100);not null;index" json:"slug"`
	Rank        uint                  `gorm:"default:0" json:"rank"`
	IsPublic    *utilities.NullBool   `gorm:"type:bool;default:true" json:"is_public"`
	IsActive    *utilities.NullBool   `gorm:"type:bool;default:true" json:"is_active"`
	CreatedAt   time.Time             `gorm:"type:datetime" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"type:datetime" json:"updated_at"`
}

func DropRoleTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&Role{})
	if err != nil {
		return err
	}

	fmt.Println("Drop Role Table")
	return nil
}
