package login

import (
	"backend/src/authentication"
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
	loginRepo 	interfaces.LoginRepositoryInterface
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
		loginRepo = repositories.NewPostgreUserRepository()
	})
}

func GetLoginRepository() (interfaces.LoginRepositoryInterface, error) {
	if loginRepo != nil {
		return loginRepo, repoErr
	}
	
	initRepository()
	if repoErr != nil {
		return nil, repoErr
	}

	return loginRepo, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		exceptions.HandleError(w, r, http.StatusUnprocessableEntity, exceptions.ErrBadRequest)
		return
	}

	var loginData model.LoginUser
	if err := json.Unmarshal(bodyRequest, &loginData); err != nil {
		exceptions.HandleError(w, r, http.StatusBadRequest, exceptions.ErrBadRequest)
		return
	}

	repo, err := GetLoginRepository()
	if err != nil {
		log.Printf("Error getting login repository: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	user, err := repo.GetUserByEmail(loginData.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		exceptions.HandleError(w ,r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	if err := user.CheckPassword(user.Password, loginData.Password); err != nil {
		exceptions.HandleError(w, r, http.StatusUnauthorized, exceptions.ErrInvalidCredentials)
		return
	}

	token, err := authentication.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		exceptions.HandleError(w, r, http.StatusInternalServerError, exceptions.ErrInternalServer)
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}