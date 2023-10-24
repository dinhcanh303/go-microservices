package domain

import (
	"time"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
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

type PostExtra struct {
	ID        uuid.UUID                          `json:"id"`
	Status    int32                              `json:"status"`
	Title     string                             `json:"title"`
	Content   string                             `json:"content"`
	UserID    uuid.UUID                          `json:"user_id"`
	GroupID   uuid.NullUUID                      `json:"group_id"`
	CreatedAt time.Time                          `json:"created_at"`
	UpdatedAt time.Time                          `json:"updated_at"`
	Likes     []*domain.Like                     `json:"likes"`
	Comments  []*sharedkernel.CommentHasChildren `json:"comments"`
}
