package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type DTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) IsValid() error {
	if u.Username == "" {
		return errors.New("username can't be empty")
	}
	if u.Email == "" {
		return errors.New("email can't be empty")
	}
	return nil
}
