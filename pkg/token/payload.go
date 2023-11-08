package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	AvatarUrl string    `json:"avatar_url"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid implements jwt.Claims.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(id uuid.UUID, email, fullName, role, avatarUrl string, duration time.Duration) *Payload {
	return &Payload{
		ID:        id,
		Email:     email,
		FullName:  fullName,
		Role:      role,
		AvatarUrl: avatarUrl,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}
