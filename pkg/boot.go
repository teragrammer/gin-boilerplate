package pkg

import (
	"fmt"
	"gin-boilerplate/configs"
	"gin-boilerplate/database"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"path/filepath"
	"strconv"
)

func InitBoot(envPath string, overrideEnv *string) configs.BootHandlers {
	// prepare the configuration file
	absEnvPath, err := filepath.Abs(envPath)
	fmt.Println("path.data", absEnvPath, envPath)
	if err != nil {
		log.Println("Error during Unmarshal(): ", err)
	}
	configurations := configs.Config(absEnvPath)
	if overrideEnv != nil {
		configurations.Environment = *overrideEnv
	}

	engine := gin.Default()

	// connect to MySQL
	db := database.Connector(configurations.Db)

	// connect to Redis
	rd := redis.NewClient(&redis.Options{
		Addr:     configurations.Redis.Host + ":" + strconv.Itoa(configurations.Redis.Port),
		Password: configurations.Redis.Password,
	})

	// initialized all configs
	var bootstrap configs.BootHandlers
	bootstrap.Env = configurations
	bootstrap.Engine = engine
	bootstrap.DB = db
	bootstrap.Redis = rd
	return bootstrap
}
