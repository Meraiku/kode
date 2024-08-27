package database

import (
	"time"

	"github.com/google/uuid"
)

func (db *DB) CreateNote(userId, body, title string) error {
	n := Note{
		ID:        uuid.New(),
		Title:     title,
		Body:      body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserId:    userId,
	}

	query := `INSERT INTO notes (id, title, body, created_at, updated_at, user_id)
	VALUES ($1 $2 $3 $4 $5 $6)`

	if _, err := db.db.Exec(query,
		&n.ID,
		&n.Title,
		&n.Body,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
	); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserNotes(userId string) ([]Note, error) {
	notes := []Note{}

	query := `SELECT * FROM notes WHERE user_id=$1 ORDER BY created_at DESC`

	rows, err := db.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := Note{}

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Body,
			&n.CreatedAt,
			&n.UpdatedAt,
			&n.UserId,
		); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}
