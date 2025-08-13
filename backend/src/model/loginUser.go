package model

import "backend/src/security"

type LoginUser struct {
	ID		 	uint64 `json:"id,omitempty"`
	Email		string `json:"email"`
	Password	string `json:"password"`
}

func (l *LoginUser) CheckPassword(passwordHashed string, password string) error {
	if err := security.VerifyPassword(passwordHashed, password); err != nil {
		return &ValidationError{
			Field:   "password",
			Message: "Invalid password",
			Code:    ErrCodeInvalidCredentials,
		}
	}

	return nil
}