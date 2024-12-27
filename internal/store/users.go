package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"role_id"`
	Role      Role     `json:"role"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type userStore struct {
	db *sql.DB
}

func (u *userStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, password, email, role_id)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
		user.RoleID,
	).Scan(
		&user.ID,
		&user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return nil
}

func (u *userStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at, roles.id, roles.name, roles.level, roles.description
		FROM users 
		JOIN roles ON roles.id = users.role_id
		WHERE users.id = $1 AND is_active = true 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := u.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.Id,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	user.RoleID = user.Role.Id

	return &user, nil
}

func (u *userStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at, roles.id, roles.name, roles.level, roles.description
		FROM users 
		JOIN roles ON roles.id = users.role_id
		WHERE users.email = $1 AND is_active = true 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := u.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.Id,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	user.RoleID = user.Role.Id

	return &user, nil
}

func (u *userStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTransaction(u.db, ctx, func(tx *sql.Tx) error {
		// transaction wrapper
		// create the user
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invite
		err := u.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *userStore) Delete(ctx context.Context, userID int64) error {
	return withTransaction(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (u *userStore) Activate(ctx context.Context, token string) error {
	return withTransaction(u.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user that this token belongs to
		user, err := u.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. clean the invitations
		if err := u.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})

}

func (u *userStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users AS u 
		JOIN users_invitations AS i ON i.user_id = u.id
		WHERE i.token = $1 AND i.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err

		}
	}

	return &user, nil
}

func (u *userStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, is_active = $4
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.IsActive)
	return err
}

func (u *userStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (u *userStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `INSERT INTO users_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (u *userStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}
