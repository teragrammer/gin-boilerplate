package database

import (
	"gin-boilerplate/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

func Connector(config configs.Db) *gorm.DB {
	dsn := config.User + ":" + config.Password +
		"@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" +
		config.Name + "?charset=" + config.Charset +
		"&parseTime=True&loc=UTC"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		panic("Failed to connect migration.")
	}

	return db
}
