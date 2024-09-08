package user

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll() ([]User, error) {
	var users []User
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	fmt.Println(3)
	return users, nil
}

func (r *Repository) Create(newUser User) error {
	query := "INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, newUser.Username, newUser.Email, newUser.Password, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
