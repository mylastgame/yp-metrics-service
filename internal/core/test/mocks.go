package test

import (
	"context"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
)

type MockFileStorage struct {
	repo storage.Repo
}

func NewMockFileStorage(repo storage.Repo) *MockFileStorage {
	return &MockFileStorage{repo: repo}
}

func (s *MockFileStorage) Save(ctx context.Context) error {
	return nil
}

func (s *MockFileStorage) Restore(ctx context.Context) error {
	return nil
}
