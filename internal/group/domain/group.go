package domain

import (
	"time"

	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type Group struct {
	sharedkernel.AggregateRoot
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      int32     `json:"status"`
	UserID      uuid.UUID `json:"user_id"`
	ProfileUrl  string    `json:"profile_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
