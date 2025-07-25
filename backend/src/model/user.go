package model

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID	   		uint64    	`json:"id,omitempty"`
	Username 	string		`json:"username,omitempty"`
	Nickname 	string		`json:"nickname,omitempty"`
	Email   	string		`json:"email,omitempty"`
	Password 	string		`json:"password,omitempty"`
	CreatedAt 	time.Time	`json:"created_at,omitempty"`
}

func (u *User) BeforeCreate() error {
	if err := u.validUser(); err != nil {
		return err
	}

	u.formatInput()
	return nil
}

func (u *User) validUser() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.Nickname == "" {
		return errors.New("nickname is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (u *User) formatInput() {
	u.Username = strings.TrimSpace(u.Username)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Email = strings.TrimSpace(u.Email)
}