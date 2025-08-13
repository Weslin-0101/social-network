package exceptions

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrInvalidUserNickname = errors.New("invalid user nickname")
	ErrDatabaseConnection = errors.New("failed to connect to the database")
	ErrInternalServer = errors.New("internal server error occurred")
	ErrBadRequest = errors.New("invalid request body format")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized = errors.New("unauthorized access")
)