package controllers

import (
	"backend/src/database"
	"backend/src/exceptions"
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
		err := database.ConnectDB()
		if err != nil {
			repoErr = err
			return
		}
		userRepo = repositories.NewPostgreUserRepository()
	})
}

func GetUserRepository() (interfaces.UserRepositoryInterface, error) {
	if userRepo != nil {
		return userRepo, repoErr
	}
	
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
		exceptions.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		exceptions.HandleError(w, http.StatusBadRequest, err)
		return
	}

	var user model.User
	if err := json.Unmarshal(bodyRequest, &user); err != nil {
		exceptions.HandleError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := repo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		exceptions.HandleError(w, http.StatusInternalServerError, err)
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