package domain

import (
	shared "go-microservices/internal/pkg/shared_kernel"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      shared.Status `json:"status"`
	UserID      uuid.UUID     `json:"user_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   time.Time     `json:"deleted_at"`
}
