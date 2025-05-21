package repositories

import (
	"errors"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"gorm.io/gorm"
)

type SettingKeyValueInterface struct {
	MxLogTry  int64  `gorm:"mx_log_try"`  // Max Login Tries
	LckPrd    int64  `gorm:"lck_prd"`     // Failed Login Tries Lockout Period
	TknLth    int64  `gorm:"tkn_lth"`     // Authentication Token Key Length
	TknExp    int64  `gorm:"tkn_exp"`     // Authentication Token Expiration
	TtaReq    bool   `gorm:"tfa_req"`     // TFA Required
	TtaEmlSnd string `gorm:"tfa_eml_snd"` // TFA Email Sender
	TtaEmlSbj string `gorm:"tfa_eml_sbj"` // TFA Email Subject
}

func Settings(DB *gorm.DB, slugs []string) (*SettingKeyValueInterface, error) {
	var settings []migration.Setting
	var keyValue SettingKeyValueInterface
	if err := DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(settings); i++ {
		var slug = settings[i].Slug
		var found = utilities.IsStringValueExistOnArray(slugs, &slug)
		if found {
			switch slug {
			case "mx_log_try":
				keyValue.MxLogTry = settings[i].ConvertValue().(int64)
			case "lck_prd":
				keyValue.LckPrd = settings[i].ConvertValue().(int64)
			case "tkn_lth":
				keyValue.TknLth = settings[i].ConvertValue().(int64)
			case "tkn_exp":
				keyValue.TknExp = settings[i].ConvertValue().(int64)
			case "tfa_req":
				keyValue.TtaReq = settings[i].ConvertValue().(bool)
			case "tfa_eml_snd":
				keyValue.TtaEmlSnd = settings[i].ConvertValue().(string)
			case "tfa_eml_sbj":
				keyValue.TtaEmlSbj = settings[i].ConvertValue().(string)
			default:
				return nil, errors.New("invalid slug key name")
			}
		}
	}

	return &keyValue, nil
}
