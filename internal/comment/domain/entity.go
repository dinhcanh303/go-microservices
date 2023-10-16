package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	ReplyToID       uuid.NullUUID `json:"reply_to_id"`
	Content         string        `json:"content"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
type CommentHasChild struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	ReplyToID       uuid.NullUUID `json:"reply_to_id"`
	Content         string        `json:"content"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	Children        []*Comment    `json:"children"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
