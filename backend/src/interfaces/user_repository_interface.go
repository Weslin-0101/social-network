package interfaces

import "backend/src/model"

type UserRepositoryInterface interface {
	CreateUser(user model.User) (uint64, error)
	GetAllUsers() ([]model.User, error)
	GetUserByID(userID uint64) (model.User, error)
}