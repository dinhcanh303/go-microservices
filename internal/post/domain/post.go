package domain

import (
	"time"

	"github.com/google/uuid"
)

type PostExtra struct {
	ID        uuid.UUID      `json:"id"`
	Status    int32          `json:"status"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	UserID    uuid.UUID      `json:"user_id"`
	GroupID   uuid.NullUUID  `json:"group_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Likes     []*LikeItem    `json:"likes"`
	Comments  []*CommentItem `json:"comments"`
}

type LikeItem struct {
	ID           uuid.UUID `json:"id"`
	Emoji        string    `json:"emoji"`
	LikeableType string    `json:"likeable_type"`
	LikeableID   uuid.UUID `json:"likeable_id"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type CommentItem struct {
	ID              uuid.UUID     `json:"id"`
	Content         string        `json:"content"`
	ReplyToID       uuid.NullUUID `json:"reply_to_id"`
	UserID          uuid.UUID     `json:"user_id"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	Child           []*Comment    `json:"child"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       time.Time     `json:"deleted_at"`
}
type Comment struct {
	ID              uuid.UUID     `json:"id"`
	Content         string        `json:"content"`
	ReplyToID       uuid.NullUUID `json:"reply_to_id"`
	UserID          uuid.UUID     `json:"user_id"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       time.Time     `json:"deleted_at"`
}
