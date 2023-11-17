package domain

import (
	"time"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/google/uuid"
)

type Notification struct {
	ID             int32     `json:"id"`
	NotifiableType string    `json:"notifiable_type"`
	NotifiableID   uuid.UUID `json:"notifiable_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewNotification(e event.CommentNoti)
