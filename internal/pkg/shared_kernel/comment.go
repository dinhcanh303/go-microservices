package sharedkernel

import (
	"time"

	domainComment "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type CommentHasChildren struct {
	ID              uuid.UUID                       `json:"id"`
	UserID          uuid.UUID                       `json:"user_id"`
	ReplyToID       uuid.NullUUID                   `json:"reply_to_id"`
	Content         string                          `json:"content"`
	PostID          uuid.UUID                       `json:"post_id"`
	ParentCommentID uuid.NullUUID                   `json:"parent_comment_id"`
	Likes           []*domainLike.Like              `json:"likes"`
	Attachments     []*domainUpload.Attachment      `json:"attachments"`
	Children        []*domainComment.CommentHasLike `json:"children"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
}
