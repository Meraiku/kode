package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (db *DB) CreateUser(email string, ctx context.Context) (string, error) {
	u := &User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	query := `INSERT INTO users (id, email, created_at, updated_at)
	VALUES ($1, $2, $3, $4)`

	if _, err := db.db.Exec(query,
		u.ID,
		u.Email,
		u.CreatedAt,
		u.UpdatedAt,
	); err != nil {
		return "", err
	}

	return u.ID.String(), nil
}

func (db *DB) GetUsers(ctx context.Context) ([]User, error) {

	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT 100`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		u := &User{}

		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.RefreshToken,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			} else {
				return nil, err
			}
		}

		users = append(users, *u)
	}

	return users, nil
}

func (db *DB) GetUserByID(id string, ctx context.Context) (*User, error) {

	query := `SELECT * FROM users WHERE id=$1`

	row := db.db.QueryRow(query, id)

	u := &User{}

	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.RefreshToken,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (db *DB) UpdateUserInfo(id, refreshToken string, ctx context.Context) error {

	query := `UPDATE users SET refresh_token=$1 , updated_at=$2 WHERE id=$3`

	res, err := db.db.Exec(query, refreshToken, time.Now().UTC(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		} else {
			return err
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
