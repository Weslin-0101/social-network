package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	fmt.Println("Starting the backend server...")
	fmt.Println("Loading configuration...")
	config.Load()

	fmt.Println("Connecting to the database...")
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nShutting down gracefully...")
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		os.Exit(0)
	}()
	
	db := database.GetDB()
	migrationPath, _ := filepath.Abs("./migrations")
	if err := database.RunMigrations(db, migrationPath); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("Migrations completed successfully.")

	fmt.Println("Running the backend server...")
	router := router.Generate()

	fmt.Printf("Backend server is running on port %d\n", config.APIPort)
	fmt.Println("Database connected and ready!")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), router))
}