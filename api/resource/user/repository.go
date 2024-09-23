package user

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetById(id string) (User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id)

	var user User
	if err := row.Scan(&user.Id, &user.Username, &user.Email); err != nil {
		return user, err
	}

	return user, nil
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
		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) Create(newUser User) error {
	query := "INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, newUser.Username, newUser.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Update(user User) error {
	query := "UPDATE users SET username = $1, email=$2, password=$3, updated_at=$4 WHERE id=$5"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
