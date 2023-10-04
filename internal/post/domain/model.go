package domain

import (
	"time"

	shared "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID     `json:"id"`
	Status    shared.Status `json:"status"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	UserID    string        `json:"user_id"`
	GroupID   string        `json:"group_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt time.Time     `json:"deleted_at"`
}
