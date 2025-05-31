package main

import (
	"backend/src/router"
	"backend/tests/redis"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Running the backend server...")
	router := router.Generate()

	err := redis.PopulateUsers()
	if err != nil {
		log.Fatalf("Error populating users: %v", err)
	}

	fmt.Println("Users populated successfully")
	fmt.Println("Backend server is running on port 5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}