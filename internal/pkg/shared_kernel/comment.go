package sharedkernel

import (
	"time"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/google/uuid"
)

type CommentHasChildren struct {
	ID              uuid.UUID         `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	ReplyToID       uuid.NullUUID     `json:"reply_to_id"`
	Content         string            `json:"content"`
	PostID          uuid.UUID         `json:"post_id"`
	ParentCommentID uuid.NullUUID     `json:"parent_comment_id"`
	Children        []*domain.Comment `json:"children"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
