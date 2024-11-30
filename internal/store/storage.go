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
		Create(ctx context.Context, user *User) error
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
