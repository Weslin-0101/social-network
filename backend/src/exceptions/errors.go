package exceptions

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID")
)