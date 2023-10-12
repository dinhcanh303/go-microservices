package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	ReplyTo         string        `json:"reply_to"`
	Content         string        `json:"content"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       time.Time     `json:"deleted_at"`
}
