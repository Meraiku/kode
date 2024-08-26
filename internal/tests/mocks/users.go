package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/kode/internal/database"
)

func (m *MockDB) CreateUser(email string, ctx context.Context) (string, error) {

	for _, v := range m.DB {
		if email == v.Email {
			return "", database.ErrNotFound
		}
	}

	user := database.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	m.DB[user.ID.String()] = user

	return user.ID.String(), nil
}

func (m *MockDB) GetUsers(ctx context.Context) ([]database.User, error) {

	users := []database.User{}

	for _, v := range m.DB {
		users = append(users, v)
	}

	return users, nil
}

func (m *MockDB) GetUserByID(id string, ctx context.Context) (*database.User, error) {

	user, ok := m.DB[id]
	if !ok {
		return nil, database.ErrNotFound
	}

	return &user, nil
}

func (m *MockDB) UpdateUserInfo(id, refreshToken string, ctx context.Context) error {

	user, ok := m.DB[id]
	if !ok {
		return database.ErrNotFound
	}
	user.RefreshToken = &refreshToken

	m.DB[id] = user

	return nil
}
