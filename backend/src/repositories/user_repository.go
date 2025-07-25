package repositories

import (
	"backend/src/database"
	"backend/src/model"
	"database/sql"
	"fmt"
)

type PostgreUserRepository struct {
	db *sql.DB
}

func NewPostgreUserRepository() *PostgreUserRepository {
	return &PostgreUserRepository{
		db: database.DB,
	}
}

func (r *PostgreUserRepository) CreateUser(user model.User) (uint64, error) {
	query := `
		INSERT INTO users (username, nickname, email, password) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	var userID uint64
	err := r.db.QueryRow(
		query, 
		user.Username, 
		user.Nickname, 
		user.Email,
		user.Password,
	).Scan(&userID)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}