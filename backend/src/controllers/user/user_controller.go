package controllers

import (
	"backend/src/database"
	"backend/src/interfaces"
	"backend/src/model"
	"backend/src/repositories"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	userRepo 	interfaces.UserRepositoryInterface
	repoOnce 	sync.Once
	repoErr 	error
)

func initRepository() {
	repoOnce.Do(func() {
		redisClient, err := database.ConnectRedis()
		if err != nil {
			repoErr = err
			return
		}
		userRepo = repositories.NewRedisUserRepository(redisClient)
	})
}

func GetUserRepository() (interfaces.UserRepositoryInterface, error) {
	initRepository()
	if repoErr != nil {
		return nil, repoErr
	}

	return userRepo, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.Unmarshal(bodyRequest, &user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userID, err := repo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"user_id": userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All users retrieved successfully"))
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User retrieved successfully"))
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User updated successfully"))
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted successfully"))
}