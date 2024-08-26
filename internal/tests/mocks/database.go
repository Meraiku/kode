package mocks

import (
	"github.com/meraiku/kode/internal/database"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	DB map[string]database.User
	mock.Mock
}

func NewMockDB() *MockDB {
	return &MockDB{
		DB: map[string]database.User{},
	}
}
