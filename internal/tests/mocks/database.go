package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	DB map[string]any
	mock.Mock
}

func NewMockDB() *MockDB {
	return &MockDB{
		DB: map[string]any{},
	}
}
