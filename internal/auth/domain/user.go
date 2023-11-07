package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FistName  string    `json:"fist_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email string, password string, fistName string, lastName string) *User {
	return &User{
		ID:       uuid.New(),
		Email:    email,
		FistName: fistName,
		LastName: lastName,
		FullName: fistName + " " + lastName,
		Password: password,
	}
}
