package model

import "time"

type User struct {
	ID	   		uint64    	`json:"id,omitempty"`
	Username 	string		`json:"username,omitempty"`
	Nickname 	string		`json:"nickname,omitempty"`
	Email   	string		`json:"email,omitempty"`
	Password 	string		`json:"password,omitempty"`
	CreatedAt 	time.Time	`json:"created_at"`
}