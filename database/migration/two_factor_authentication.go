package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type TwoFactorAuthentication struct {
	Id             uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenId        uint                `gorm:"not null" json:"token_id"`
	Code           string              `gorm:"varchar(100)" json:"code"`
	Tries          uint                `gorm:"type:int(10)" json:"tries"`
	NextSendAt     *utilities.NullTime `gorm:"type:TIMESTAMP" json:"next_send_at"`
	ExpiredTriesAt *utilities.NullTime `gorm:"type:TIMESTAMP" json:"expired_tries_at"`
	CreatedAt      time.Time           `gorm:"type:datetime" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"type:datetime" json:"updated_at"`
	ExpiredAt      *utilities.NullTime `gorm:"type:datetime" json:"expired_at"`

	AuthenticationToken *AuthenticationToken `gorm:"foreignKey:TokenId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"authentication_token,omitempty"`
}

func DropTwoFactorAuthenticationTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&TwoFactorAuthentication{})
	if err != nil {
		return err
	}

	fmt.Println("Drop Two Factor Authentication Table")
	return nil
}
