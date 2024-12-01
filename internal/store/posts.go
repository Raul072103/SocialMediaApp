package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	User      User       `json:"user"`
	Tags      []string   `json:"tags"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  *[]Comment `json:"comments"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}

type postStore struct {
	db *sql.DB
}

func (s *postStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	//shouldFilterByOwnPost := len(fq.Tags) >= 1 || fq.Search != ""
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, COUNT(c.id), u.username AS comments_count
		FROM posts AS p
		LEFT JOIN comments AS c ON c.post_id = p.id
		LEFT JOIN users AS u ON p.user_id = u.id
		JOIN followers AS f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
		     f.user_id = $1 AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort +
		` LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
		fq.Limit,
		fq.Offset,
		fq.Search,
		pq.Array(fq.Tags),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.CommentCount,
			&p.User.Username,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

func (s *postStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *postStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT id, title, user_id, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.UserID,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *postStore) DeleteById(ctx context.Context, postID int64) error {
	query := `
		DELETE
		FROM posts
		WHERE posts.id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *postStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
