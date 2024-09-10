package user

import (
	"errors"
	"time"
)

type User struct {
	Id        string
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
