package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type Setting struct {
	Id          uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string                `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string                `gorm:"type:varchar(100);unique;not null" json:"slug"`
	Value       *utilities.NullString `gorm:"type:TEXT" json:"value"`
	Description *utilities.NullString `gorm:"type:TINYTEXT" json:"description"`
	Type        string                `gorm:"type:ENUM('string', 'integer', 'float', 'boolean', 'array');default:string" json:"type"`
	IsDisabled  *utilities.NullBool   `gorm:"type:bool;default:false" json:"is_disabled"`
	IsPublic    *utilities.NullBool   `gorm:"type:bool;default:true" json:"is_public"`
	CreatedAt   time.Time             `gorm:"type:datetime" json:"created_at"`
	UpdatedAt   time.Time             `gorm:"type:datetime" json:"updated_at"`
}

func DropSettingTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&Setting{})
	if err != nil {
		return err
	}

	fmt.Println("Drop Setting Table")
	return nil
}
