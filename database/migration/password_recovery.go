package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type Type string

const (
	Email Type = "email"
	Phone Type = "phone"
)

type PasswordRecovery struct {
	Id           uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	Type         Type                  `gorm:"type:ENUM('email','phone')"`
	SendTo       *utilities.NullString `gorm:"type:varchar(100);unique" json:"send_to"`
	Code         string                `gorm:"type:varchar(100)" json:"code"`
	NextResendAt time.Time             `gorm:"type:datetime" json:"next_resend_at"`
	ExpiredAt    time.Time             `gorm:"type:datetime" json:"expired_at"`
	Tries        uint                  `gorm:"default:0" json:"tries"`
	NextTryAt    time.Time             `gorm:"type:datetime" json:"next_try_at"`
	CreatedAt    time.Time             `gorm:"type:datetime" json:"created_at"`
	UpdatedAt    time.Time             `gorm:"type:datetime" json:"updated_at"`
}

func DropPasswordRecoveryTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&PasswordRecovery{})
	if err != nil {
		return err
	}

	fmt.Println("Drop Password Recovery Table")
	return nil
}
