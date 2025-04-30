package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type AuthenticationToken struct {
	Id            uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId        uint                `gorm:"not null" json:"user_id"`
	Token         string              `gorm:"type:varchar(180);unique;not null" json:"token"`
	IsTFARequired *utilities.NullBool `gorm:"type:bool;default:false" json:"is_tfa_required"`
	IsTFAVerified *utilities.NullBool `gorm:"type:bool;default:false" json:"is_tfa_verified"`
	CreatedAt     time.Time           `gorm:"type:datetime" json:"created_at"`
	UpdatedAt     time.Time           `gorm:"type:datetime" json:"updated_at"`
	ExpiredAt     *utilities.NullTime `gorm:"type:datetime" json:"expired_at"`

	User *User `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func DropAuthenticationTokenTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&AuthenticationToken{})
	if err != nil {
		return err
	}

	fmt.Println("Drop Authentication Token Table")
	return nil
}
