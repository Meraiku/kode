package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/kode/internal/database"
)

func (db *MockDB) CreateNote(userId, body, title string) (database.Note, error) {
	n := database.Note{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserId:    userId,
	}

	db.DB[n.ID.String()] = n

	return n, nil
}

func (db *MockDB) GetUserNotes(userId string) ([]database.Note, error) {
	notes := []database.Note{}

	for _, vt := range db.DB {
		switch v := vt.(type) {
		case database.Note:
			if v.UserId == userId {
				notes = append(notes, v)
			}
		}
	}

	return notes, nil
}
