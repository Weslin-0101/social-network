package repositories

import (
	"backend/src/database"
	"backend/src/exceptions"
	"backend/src/model"
	"database/sql"
	"fmt"
	"strings"
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

func (r *PostgreUserRepository) GetUserByID(userID uint64) (model.User, error) {
	query := `
		SELECT
			id, username, nickname, email, created_at
		FROM
			users
		WHERE
			id = $1
	`

	var user model.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, exceptions.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *PostgreUserRepository) GetUserByNickname(nickname string) (model.User, error) {
	query := `
		SELECT
			id, username, nickname, email, created_at
		FROM
			users
		WHERE
			nickname = $1
	`

	var user model.User
	err := r.db.QueryRow(query, strings.TrimSpace(nickname)).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, exceptions.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("failed to get user by nickname: %w", err)
	}

	return user, nil
}