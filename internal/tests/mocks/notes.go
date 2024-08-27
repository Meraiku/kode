package mocks

import "github.com/meraiku/kode/internal/database"

func (db *MockDB) CreateNote(userId, body, title string) error {
	return nil
}

func (db *MockDB) GetUserNotes(userId string) ([]database.Note, error) {
	return nil, nil
}
