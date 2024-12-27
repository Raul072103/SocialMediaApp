package store

import (
	"context"
	"database/sql"
)

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}

type rolesStore struct {
	db *sql.DB
}

func (r *rolesStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `SELECT id, name, level, description
	FROM roles
	WHERE name = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var role Role

	err := r.db.QueryRowContext(
		ctx,
		query,
		roleName,
	).Scan(
		&role.Id,
		&role.Name,
		&role.Level,
		&role.Description)

	if err != nil {
		return nil, err
	}

	return &role, nil
}
