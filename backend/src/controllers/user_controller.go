package controllers

import "net/http"

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User created successfully"))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All users retrieved successfully"))
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User retrieved successfully"))
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User updated successfully"))
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted successfully"))
}