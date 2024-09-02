package user

import (
	"time"

	"github.com/google/uuid"
)

type DTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
