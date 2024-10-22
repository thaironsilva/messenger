package user

import (
	"database/sql"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByUsername(username string) (User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE username = $1", username)

	var user User
	if err := row.Scan(&user.Id, &user.Username, &user.Email); err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(email string) (User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email)

	var user User
	if err := row.Scan(&user.Id, &user.Username, &user.Email); err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) GetByString(name string) ([]User, error) {
	name = "%" + strings.ToLower(name) + "%"
	rows, err := r.db.Query("SELECT * FROM users WHERE username LIKE $1 OR email LIKE $1 LIMIT 20", name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

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
	query := "INSERT INTO users (username, email) VALUES ($1, $2)"
	_, err := r.db.Exec(query, newUser.Username, newUser.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Update(user User) error {
	query := "UPDATE users SET username = $1 WHERE id=$5"
	_, err := r.db.Exec(query, user.Username, user.Id)
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
