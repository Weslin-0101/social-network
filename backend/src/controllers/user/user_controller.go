package controllers

import (
	"backend/src/authentication"
	"backend/src/database"
	"backend/src/exceptions"
	"backend/src/interfaces"
	"backend/src/model"
	"backend/src/repositories"
	"encoding/json"
	"fmt"
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
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrBadRequest)
		return
	}

	var user model.User
	if err := json.Unmarshal(bodyRequest, &user); err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrBadRequest)
		return
	}

	if err := user.BeforeCreate("register"); err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, err)
		return
	}

	user, err = repo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	users, err := repo.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
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
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrInvalidUserID)
		return
	}

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	user, err := repo.GetUserByID(uint64(userID))
	if err != nil {
		if err == exceptions.ErrUserNotFound {
			log.Printf("User with ID %d not found", userID)
			exceptions.HandleError(w, r, http.StatusNotFound, exceptions.ErrUserNotFound)
			return
		}

		log.Printf("Error retrieving user by ID: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	if user == (model.User{}) {
		exceptions.HandleError(w, r, http.StatusNotFound, nil)
		return
	}

	response := map[string]interface{}{
		"message": 	"User retrieved successfully",
		"user": 	user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]
	if nickname == "" {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrInvalidUserNickname)
		return
	}

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	user, err := repo.GetUserByNickname(nickname)
	if err != nil {
		if err == exceptions.ErrUserNotFound {
			log.Printf("User with nickname %s not found", nickname)
			exceptions.HandleError(w, r, http.StatusNotFound, exceptions.ErrUserNotFound)
			return
		}
		log.Printf("Error retrieving user by nickname: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	response := map[string]interface{} {
		"message": "User retrieved successfully",
		"user": 	user,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrInvalidUserID)
		return
	}

	userIDFromToken, err := authentication.ExtractUserID(r)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusUnauthorized, exceptions.ErrUnauthorized)
		return
	}

	if userID != userIDFromToken {
		exceptions.HandleError(w, r, http.StatusForbidden, exceptions.ErrForbidden)
		return
	}

	fmt.Println(userIDFromToken)

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusUnprocessableEntity, exceptions.ErrBadRequest)
		return
	}

	var user model.User
	if err := json.Unmarshal(bodyRequest, &user); err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrBadRequest)
		return
	}

	if err = user.BeforeCreate("update"); err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrBadRequest)
		return
	}

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrDatabaseConnection)
		return
	}

	_, err = repo.UpdateUserByID(userID, user)
	if err != nil {
		if err == exceptions.ErrUserNotFound {
			log.Printf("User with ID %d not found", userID)
			exceptions.HandleError(w, r, http.StatusNotFound, exceptions.ErrUserNotFound)
			return
		}
		log.Printf("Error updating user by ID: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrInvalidUserID)
		return
	}

	repo, err := GetUserRepository()
	if err != nil {
		log.Printf("Error getting user repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrDatabaseConnection)
		return
	}

	err = repo.DeleteUserByID(userID)
	if err != nil {
		if err == exceptions.ErrUserNotFound {
			log.Printf("User with ID %d not found", userID)
			exceptions.HandleError(w, r, http.StatusNotFound, exceptions.ErrUserNotFound)
			return
		}
		log.Printf("Error deleting user by ID: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrDatabaseConnection)
		return
	}

	response := map[string]string{
		"message": "User deleted successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}