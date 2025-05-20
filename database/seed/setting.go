package seed

import (
	"database/sql"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
)

func SeedingSettingTable(db *gorm.DB) error {
	settings := []migration.Setting{
		{Name: "Max Login Tries", Slug: "mx_log_try", Value: &utilities.NullString{NullString: sql.NullString{String: "5", Valid: true}}, Type: "integer", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "Failed Login Tries Lockout Period", Slug: "lck_prd", Value: &utilities.NullString{NullString: sql.NullString{String: "5", Valid: true}}, Type: "integer", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "Authentication Token Key Length", Slug: "tkn_lth", Value: &utilities.NullString{NullString: sql.NullString{String: "62", Valid: true}}, Type: "integer", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "Authentication Token Expiration", Slug: "tkn_exp", Value: &utilities.NullString{NullString: sql.NullString{String: "15", Valid: true}}, Type: "integer", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "TFA Required", Slug: "tfa_req", Value: &utilities.NullString{NullString: sql.NullString{String: "0", Valid: true}}, Type: "boolean", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "TFA Email Sender", Slug: "tfa_eml_snd", Value: nil, Type: "string", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		{Name: "TFA Email Subject", Slug: "tfa_eml_sbj", Value: &utilities.NullString{NullString: sql.NullString{String: "Your OTP Code", Valid: true}}, Type: "string", IsPublic: &utilities.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
	}

	// Insert records into the database
	if err := db.Create(&settings).Error; err != nil {
		return err
	} else {
		return nil
	}
}
