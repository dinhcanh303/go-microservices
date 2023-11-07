package domain

import "time"

type ApiKey struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
