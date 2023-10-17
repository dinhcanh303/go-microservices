package domain

import (
	"time"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type PostExtra struct {
	ID        uuid.UUID                          `json:"id"`
	Status    int32                              `json:"status"`
	Title     string                             `json:"title"`
	Content   string                             `json:"content"`
	UserID    uuid.UUID                          `json:"user_id"`
	GroupID   uuid.NullUUID                      `json:"group_id"`
	CreatedAt time.Time                          `json:"created_at"`
	UpdatedAt time.Time                          `json:"updated_at"`
	Likes     []*LikeItem                        `json:"likes"`
	Comments  []*sharedkernel.CommentHasChildren `json:"comments"`
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
