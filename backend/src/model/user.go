package model

import (
	"backend/src/security"
	"regexp"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID	   		uint64    	`json:"id,omitempty"`
	Username 	string		`json:"username,omitempty"`
	Nickname 	string		`json:"nickname,omitempty"`
	Email   	string		`json:"email,omitempty"`
	Password 	string		`json:"password,omitempty"`
	Type		string		`json:"type,omitempty"`
	CreatedAt 	time.Time	`json:"created_at,omitempty"`
}

type ValidationError struct {
	Field  	string `json:"field"`
	Message string `json:"message"`
	Code  	string `json:"code"`
}

func (ve ValidationError) Error() string {
	return ve.Message
}

const (
	ErrCodeRequired			= "FIELD_REQUIRED"
	ErrCodeInvalidFormat 	= "INVALID_FORMAT"
	ErrCodeTooShort			= "TOO_SHORT"
	ErrCodeTooLong			= "TOO_LONG"
	ErrCodeInvalidChars 	= "INVALID_CHARACTERS"
	ErrCodeTooWeakPassword 	= "PASSWORD_TOO_WEAK"
)

func (u *User) BeforeCreate(step string) error {
	if err := u.validUser(step); err != nil {
		return err
	}

	if err := u.formatInput(step); err != nil {
		return err
	}
	
	return nil
}

func (u *User) validUser(step string) error {
	if err := u.validUsername(); err != nil { return err }

	if err := u.validNickname(); err != nil { return err }

	if err := u.validEmail(); err != nil { return err }

	if step == "register" {
		if err := u.validPassword(); err != nil { return err }
	}

	return nil
}

func (u *User) validUsername() error {
	username := strings.TrimSpace(u.Username)

	if username == "" {
		return ValidationError {
			Field:   "username",
			Message: "Username is required",
			Code:    ErrCodeRequired,
		}
	}

	if len(username) < 3 {
		return ValidationError {
			Field:   "username",
			Message: "Username must be at least 3 characters long",
			Code:    ErrCodeTooShort,
		}
	}

	if len(username) > 100 {
		return ValidationError {
			Field:   "username",
			Message: "Username must be at most 100 characters long",
			Code:    ErrCodeTooLong,
		}
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return ValidationError {
			Field:  	"username",
			Message: 	"Username can only contain letters, numbers, underscores, and hyphens",
			Code:   	ErrCodeInvalidChars,
		}
	}

	return nil
}

func (u *User) validNickname() error {
	nickname := strings.TrimSpace(u.Nickname)

	if nickname == "" {
		return ValidationError {
			Field:   "nickname",
			Message: "Nickname is required",
			Code:    ErrCodeRequired,
		}
	}

	if len(nickname) < 3 {
		return ValidationError {
			Field:   "nickname",
			Message: "Nickname must be at least 3 characters long",
			Code:    ErrCodeTooShort,
		}
	}

	if len(nickname) > 100 {
		return ValidationError {
			Field:   "nickname",
			Message: "Nickname must be at most 100 characters long",
			Code:    ErrCodeTooLong,
		}
	}

	return nil
}

func (u *User) validEmail() error {
	email := strings.TrimSpace(u.Email)

	if email == "" {
		return ValidationError {
			Field:   "email",
			Message: "Email is required",
			Code:    ErrCodeRequired,
		}
	}

	if len(email) > 100 {
		return ValidationError{
			Field:   "email",
			Message: "Email must be at most 100 characters long",
			Code:    ErrCodeTooLong,
		}
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		return ValidationError {
			Field:   "email",
			Message: "Invalid email format",
			Code:    ErrCodeInvalidFormat,
		}
	}

	return nil
}

func (u *User) validPassword() error {
	if u.Password == "" {
		return ValidationError {
			Field:   "password",
			Message: "Password is required",
			Code:    ErrCodeRequired,
		}
	}

	if len(u.Password) < 8 {
		return ValidationError {
			Field:   "password",
			Message: "Password must be at least 8 characters long",
			Code:    ErrCodeTooShort,
		}
	}

	if len(u.Password) > 100 {
		return ValidationError {
			Field:   "password",
			Message: "Password must be at most 100 characters long",
			Code:    ErrCodeTooLong,
		}
	}

	if security.IsStrongPassword(u.Password) {
		return ValidationError {
			Field:   "password",
			Message: "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
			Code:    ErrCodeTooWeakPassword,
		}
	}

	return nil
}

func (u *User) formatInput(step string) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Email = strings.TrimSpace(u.Email)

	if step == "register" {
		passwordHash, err := security.HashPassword(u.Password)
		if err != nil {
			return ValidationError {
				Field:		"password",
				Message: 	"Failed to process password",
				Code:		"HASH_ERROR",
			}
		}

		u.Password = string(passwordHash)
	}

	return nil
}