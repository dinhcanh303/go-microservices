package domain

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `json:"id"`
	Emoji     string    `json:"emoji"`
	UserID    uint64    `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	CommentID uuid.UUID `json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
