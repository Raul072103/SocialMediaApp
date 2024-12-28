package cache

import (
	"SocialMediaApp/internal/store"
	"context"
)

func NewMockStorage() Storage {
	return Storage{Users: &MockUserStore{}}
}

type MockUserStore struct{}

func (m *MockUserStore) Get(context.Context, int64) (*store.User, error) {
	return nil, nil
}

func (m *MockUserStore) Set(context.Context, *store.User) error {
	return nil
}
