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
	"strconv"
	"sync"

	"github.com/gorilla/mux"
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

	if err := user.BeforeCreate(); err != nil {
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
	w.Header().Set("Content-Type", "application/json")

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	users, err := repo.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		exceptions.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		exceptions.HandleError(w, http.StatusBadRequest, err)
		return
	}

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := repo.GetUserByID(uint64(userID))
	if err != nil {
		if err == exceptions.ErrUserNotFound {
			log.Printf("User with ID %d not found", userID)
			exceptions.HandleError(w, http.StatusNotFound, err)
			return
		}

		log.Printf("Error retrieving user by ID: %v", err)
		exceptions.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	if user == (model.User{}) {
		exceptions.HandleError(w, http.StatusNotFound, nil)
		return
	}

	response := map[string]interface{}{
		"message": 	"User retrieved successfully",
		"user": 	user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User updated successfully"))
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted successfully"))
}