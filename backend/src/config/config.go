package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	RedisHost = ""
	RedisPort = 0
	RedisDB   = 0
	ApiPort   = 0
)

func Load() {
	var err error

	// Try to load environment variables from a .env file but don't fail if it doesn't exist
	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values from system")
	}

	RedisPort, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		RedisPort = 6379
	}

	RedisDB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		RedisDB = 0
	}

	RedisHost = fmt.Sprintf("%s:%d", os.Getenv("REDIS_HOST"), RedisPort)

	ApiPort, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		ApiPort = 5000
	}
}