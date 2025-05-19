package migration

import (
	"fmt"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id                   uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName            string                `gorm:"type:varchar(100);not null;index" json:"first_name"`
	MiddleName           *utilities.NullString `gorm:"type:varchar(100);" json:"middle_name"`
	LastName             *utilities.NullString `gorm:"type:varchar(100);index;" json:"last_name"`
	Gender               *utilities.NullString `gorm:"type:ENUM('Male', 'Female');" json:"gender"`
	Address              *utilities.NullString `gorm:"type:TINYTEXT" json:"address"`
	BirthDate            *utilities.NullTime   `gorm:"type:date" json:"birth_date"`
	Phone                *utilities.NullString `gorm:"type:varchar(50);index;unique" json:"phone"`
	IsPhoneVerified      *utilities.NullBool   `gorm:"type:bool;default:false" json:"-"`
	Email                *utilities.NullString `gorm:"type:varchar(256);index;unique" json:"email"`
	IsEmailVerified      *utilities.NullBool   `gorm:"type:bool;default:false" json:"-"`
	RoleId               uint                  `gorm:"not null" json:"role_id"`
	Username             string                `gorm:"type:varchar(50);index;unique" json:"username"`
	Password             string                `gorm:"type:TINYTEXT" json:"-"`
	Status               string                `gorm:"type:ENUM('Pending', 'Activated', 'Suspended', 'Deleted', 'Deactivated');default:Activated" json:"status"`
	LoginTries           uint                  `gorm:"type:uint;default:0" json:"-"`
	FailedLoginExpiredAt *utilities.NullTime   `gorm:"type:datetime" json:"-"`
	Comments             *utilities.NullString `gorm:"type:TINYTEXT" json:"comments"`
	CreatedAt            time.Time             `gorm:"type:datetime" gorm:"index" json:"created_at"`
	UpdatedAt            time.Time             `gorm:"type:datetime" json:"updated_at"`
	DeletedAt            *utilities.NullTime   `gorm:"type:datetime" json:"deleted_at"`

	Role *Role `gorm:"foreignKey:RoleId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"role,omitempty"`
}

func DropUserTable(db *gorm.DB) error {
	// Drop all tables
	err := db.Migrator().DropTable(&User{})
	if err != nil {
		return err
	}

	fmt.Println("Drop User Table")
	return nil
}
