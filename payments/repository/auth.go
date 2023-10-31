package repository

import (
	"context"
	"database/sql"
	"log"
	"payments/models"

	"github.com/zeebo/errs"
)

var authErr = errs.Class("authorization repository")

type AuthorizationDB struct {
	db *sql.DB
}

func NewAuthorizationDB(db *sql.DB) *AuthorizationDB {
	return &AuthorizationDB{db: db}
}

func (r *AuthorizationDB) CreateUser(ctx context.Context, user models.User) (id int, err error) {
	statement := "INSERT INTO users(name, email, password_hash) VALUES ($1, $2, $3) RETURNING id;"

	stmt, err := r.db.PrepareContext(ctx, statement)
	if err != nil {
		return 0, authErr.Wrap(err)
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, user.Name, user.Email, user.PasswordHash).Scan(&id); err != nil {
		return 0, authErr.Wrap(err)
	}

	log.Printf("Repository: You've created user:\n%#v\nwith id = %v\n", user, id)

	return id, nil
}

func (r *AuthorizationDB) GetUserID(ctx context.Context, email, passwordHash string) (id int, err error) {
	query := "SELECT id FROM users WHERE email=$1 AND password_hash=$2"
	row := r.db.QueryRowContext(ctx, query, email, passwordHash)
	if err := row.Scan(&id); err != nil {
		return 0, authErr.Wrap(err)
	}
	return id, nil
}
