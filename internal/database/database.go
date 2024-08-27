package database

import (
	"context"
	"database/sql"
)

type Store interface {
	CreateUser(string, context.Context) (string, error)
	GetUsers(context.Context) ([]User, error)
	GetUserByID(string, context.Context) (*User, error)
	UpdateUserInfo(string, string, context.Context) error
	GetUserNotes(userId string) ([]Note, error)
	CreateNote(userId, body, title string) (Note, error)
}

type DB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) Store {
	return &DB{
		db: db,
	}
}
