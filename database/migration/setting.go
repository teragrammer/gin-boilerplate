package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"strconv"
	"strings"
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

func (s *Setting) ConvertValue() interface{} {
	switch strings.ToLower(s.Type) {
	case "string":
		if s.Value == nil {
			return nil
		}
		if s.Value.Valid == false {
			return nil
		}
		return s.Value.String
	case "integer":
		if s.Value == nil {
			return int64(0)
		}
		val, _ := strconv.ParseInt(s.Value.String, 10, 64)
		return val
	case "float":
		if s.Value == nil {
			return 0.0
		}
		val, _ := strconv.ParseFloat(s.Value.String, 64)
		return val
	case "boolean":
		if s.Value == nil {
			return false
		}
		val, _ := strconv.ParseBool(s.Value.String)
		return val
	case "array":
		if s.Value == nil {
			return []string{}
		}
		if s.Value.String == "" {
			return []string{}
		}
		// Assuming array values are comma-separated strings
		vals := strings.Split(s.Value.String, ",")
		return vals
	default:
		// Handle unknown type or default case
		return nil
	}
}
