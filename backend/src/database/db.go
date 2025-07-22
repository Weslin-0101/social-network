package database

import (
	"backend/src/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	var err error

	DB, err = sql.Open("postgres", config.DBURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database connection established successfully")

	return nil
}

func GetDB() *sql.DB {
	return DB
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}

	return nil
}