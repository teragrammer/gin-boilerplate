package configs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Configuration struct {
	Env  string
	Data Env
}

type Env struct {
	Environment string
	App         App
	FileSystem  FileSystem `json:"file_system"`
	Db          Db
	Redis       Redis
	Rate        Rate
	Security    Security
}

type App struct {
	Port int
	Key  string
}

type FileSystem struct {
	Storage string
}

type Db struct {
	Host       string
	Port       int
	Name       string
	User       string
	Password   string
	PoolMin    int
	PoolMax    int
	Charset    string
	DateString bool
}

type Redis struct {
	Host     string
	Port     int
	Password string
}

type Rate struct {
	Window time.Duration
	Limit  uint
}

type Security struct {
	HashSecret    string `json:"hash_secret"`
	AESPassphrase string `json:"aes_passphrase"`
	AESSalt       string `json:"aes_salt"`
}

type BootHandlers struct {
	Env    Env
	Engine *gin.Engine
	DB     *gorm.DB
	Redis  *redis.Client
}

func Config(envPath string) Env {
	content, err := os.ReadFile(envPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	var configPayload Configuration
	err = json.Unmarshal(content, &configPayload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// set debugging mode
	if payload["env"].(string) == "development" {
		gin.SetMode(gin.DebugMode)
	}

	return configPayload.Data
}
