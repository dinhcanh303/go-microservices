package event

import (
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type CommentNoti struct {
	sharedkernel.DomainEvent
	CommentID      uuid.UUID `json:"commentId"`
	CommentContent string    `json:"commentContent"`
	CommentContent string    `json:"commentContent"`
}

// ID              uuid.UUID     `json:"id"`
// 	UserID          uuid.UUID     `json:"user_id"`
// 	ReplyToID       uuid.NullUUID `json:"reply_to_id"`
// 	Content         string        `json:"content"`
// 	PostID          uuid.UUID     `json:"post_id"`
// 	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
// 	CreatedAt       time.Time     `json:"created_at"`
// 	UpdatedAt       time.Time     `json:"updated_at"`
