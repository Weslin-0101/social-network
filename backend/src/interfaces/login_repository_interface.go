package interfaces

import "backend/src/model"

type LoginRepositoryInterface interface {
	GetUserByEmail(email string) (model.LoginUser, error)
}