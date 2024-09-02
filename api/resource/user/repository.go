package user

import (
	"database/sql"

	"github.com/google/uuid"
)

var users = []User{
	{ID: uuid.New(), Username: "user.1", Email: "user.1@email.com", Password: "password1"},
	{ID: uuid.New(), Username: "user.2", Email: "user.2@email.com", Password: "password2"},
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) list() []User {
	return users
}

func (r *Repository) create(newUser User) User {
	users = append(users, newUser)
	return newUser
}
