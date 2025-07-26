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

func (r *PostgreUserRepository) GetAllUsers() ([]model.User, error) {
	query := `
		SELECT 
			id, username, nickname, email, created_at
		FROM
			users
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	defer rows.Close()
	var users []model.User

	for rows.Next() {
		var user model.User
		if err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}