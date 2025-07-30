package interfaces

import "backend/src/model"

type UserRepositoryInterface interface {
	CreateUser(user model.User) (model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByID(userID uint64) (model.User, error)
	GetUserByNickname(nickname string) (model.User, error)
	UpdateUserByID(userID uint64, user model.User) (model.User, error)
	DeleteUserByID(userID uint64) error
}