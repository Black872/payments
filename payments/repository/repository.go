package repository

import (
	"context"
	"database/sql"
	"payments/models"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	GetUserID(ctx context.Context, name, password string) (id int, err error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthorizationDB(db),
	}
}
