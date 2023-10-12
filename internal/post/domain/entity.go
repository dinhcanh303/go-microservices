package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID     `json:"id"`
	Status    int32         `json:"status"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	UserID    uuid.UUID     `json:"user_id"`
	GroupID   uuid.NullUUID `json:"group_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
