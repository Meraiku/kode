package mocks

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/kode/internal/database"
)

func (m *MockDB) CreateUser(email string, ctx context.Context) (string, error) {

	for _, vt := range m.DB {
		switch v := vt.(type) {
		case database.User:
			if email == v.Email {
				return "", database.ErrNotFound
			}
		default:
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

	for _, vt := range m.DB {
		switch v := vt.(type) {
		case database.User:
			users = append(users, v)
		}
	}

	return users, nil
}

func (m *MockDB) GetUserByID(id string, ctx context.Context) (*database.User, error) {

	user, ok := m.DB[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	switch v := user.(type) {
	case database.User:
		return &v, nil
	default:
		return nil, database.ErrNotFound
	}

}

func (m *MockDB) UpdateUserInfo(id, refreshToken string, ctx context.Context) error {

	user, ok := m.DB[id]
	if !ok {
		return database.ErrNotFound
	}
	switch v := user.(type) {
	case database.User:
		v.RefreshToken = &refreshToken
		log.Print(*v.RefreshToken)
		m.DB[id] = v
	}

	log.Print(m.DB[id])
	return nil
}
