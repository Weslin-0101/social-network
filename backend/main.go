package main

import (
	"backend/src/config"
	"backend/src/router"
	"backend/tests/redis"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting the backend server...")
	fmt.Println("Loading configuration...")
	config.Load()
	fmt.Println("Running the backend server...")
	router := router.Generate()

	err := redis.PopulateUsers()
	if err != nil {
		log.Fatalf("Error populating users: %v", err)
	}

	fmt.Println("Users populated successfully")
	fmt.Printf("Backend server is running on port %d\n", config.ApiPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ApiPort), router))
}