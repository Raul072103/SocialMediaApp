package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("a user with this email already exists")
	ErrDuplicateUsername = errors.New("a user with this username already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		DeleteById(ctx context.Context, postID int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		GetByID(ctx context.Context, userID int64) (*User, error)
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(ctx context.Context, token string) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(context.Context, int64) (*[]Comment, error)
		DeleteByPostID(context.Context, int64) error
	}
	Followers interface {
		Follow(ctx context.Context, followedID int64, userID int64) error
		Unfollow(ctx context.Context, followedID int64, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &postStore{db},
		Users:     &userStore{db},
		Comments:  &commentStore{db},
		Followers: &followerStore{db},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
