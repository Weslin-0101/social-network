package main

import (
	"backend/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Running the backend server...")
	router := router.Generate()

	log.Fatal(http.ListenAndServe(":5000", router))
}