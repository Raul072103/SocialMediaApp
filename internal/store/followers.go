package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type followerStore struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (f *followerStore) Follow(ctx context.Context, followedID, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return err
}

func (f *followerStore) Unfollow(ctx context.Context, followedID, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 and follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, userID, followedID)
	return err
}
