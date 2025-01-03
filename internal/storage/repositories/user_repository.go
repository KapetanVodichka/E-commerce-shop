package repositories

import (
	"database/sql"
	"e-commerce-shop/internal/model/user"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u *user.User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.db.QueryRow(query, u.Username, u.Password).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*user.User, error) {
	query := `
		SELECT id, username, password
		FROM users
		WHERE username = $1
	`
	var u user.User
	err := r.db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername: %w", err)
	}
	return &u, nil
}
