package repositories

import (
	"database/sql"
	"errors"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/utilities"
	"github.com/google/uuid"
	"strconv"
	"time"
)
import "gorm.io/gorm"

func GenerateToken(DB *gorm.DB, settings *SettingKeyValueInterface, userId uint) (*migration.AuthenticationToken, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.New("Oops something went wrong during id creation, " + err.Error())
	}

	randomString, err := utilities.GenerateRandomString(int(settings.TknLth)) // token generation based on settings length
	if err != nil {
		return nil, errors.New("Oops something went wrong during token creation, " + err.Error())
	}

	secretToken := randomString + "-" + id.String()
	expiredOn := utilities.AddDay(time.Now(), int(settings.TknExp)) // expiration on how many days

	// save the token
	token := migration.AuthenticationToken{
		UserId:        userId,
		Token:         secretToken,
		IsTFARequired: &utilities.NullBool{NullBool: sql.NullBool{Valid: true, Bool: settings.TtaReq}},
		ExpiredAt:     &utilities.NullTime{NullTime: sql.NullTime{Valid: true, Time: expiredOn}},
	}
	if err := DB.Create(&token).Error; err != nil {
		return nil, errors.New("Oops something went wrong while saving token, " + err.Error())
	}

	// delete all expired token
	if err = DB.Where("user_id = ?", userId).Where("expired_at < ?", time.Now()).Delete(&migration.AuthenticationToken{}).Error; err != nil {
		return nil, errors.New("Oops something went wrong during token reset, " + err.Error())
	}

	// encode the token
	// user_id.token_id.secret.expiration
	encodedToken := utilities.EncodeBase64URL([]byte(
		strconv.Itoa(int(userId)) + "." +
			strconv.Itoa(int(token.Id)) + "." +
			secretToken + "." +
			strconv.FormatInt(expiredOn.Unix(), 10),
	))

	token.Token = encodedToken
	return &token, nil
}
