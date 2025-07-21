package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBHost 		= ""
	DBPort 		= 0
	DBUser 		= ""
	DBPassword 	= ""
	DBName 		= ""
	DBSSLMode 	= ""
	DBURL 		= ""

	APIPort 	= 0
)

func Load() {
	var err error

	// Try to load environment variables from a .env file but don't fail if it doesn't exist
	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values from system")
	}

	DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		DBPort = 5432
	}

	DBUser = os.Getenv("DB_USER")
	if DBUser == "" {
		DBUser = "postgres"
	}

	DBPassword = os.Getenv("DB_PASSWORD")
	if DBPassword == "" {
		DBPassword = "password"
	}

	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		DBName = "social_network"
	}

	DBSSLMode = os.Getenv("DB_SSLMODE")
	if DBSSLMode == "" {
		DBSSLMode = "disable"
	}

	DBURL = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		DBHost, DBPort, DBUser, DBPassword, DBName, DBSSLMode)

	APIPort, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		APIPort = 5000
	}
}