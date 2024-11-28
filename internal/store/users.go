package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type userStore struct {
	db *sql.DB
}

func (u *userStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (u *userStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users 
		WHERE id = $1
	`

	var user User

	err := u.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
