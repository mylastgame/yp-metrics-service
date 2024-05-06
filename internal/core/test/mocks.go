package test

import "github.com/mylastgame/yp-metrics-service/internal/server/storage"

type MockFileStorage struct {
	repo storage.Repo
}

func NewMockFileStorage(repo storage.Repo) *MockFileStorage {
	return &MockFileStorage{repo: repo}
}

func (s *MockFileStorage) Save() error {
	return nil
}

func (s *MockFileStorage) Restore() error {
	return nil
}
