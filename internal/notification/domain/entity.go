package domain

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         int64       `json:"id"`
	ActorID    []uuid.UUID `json:"actor_id"`
	ReceiverID []uuid.UUID `json:"receiver_id"`
	ObjectID   uuid.UUID   `json:"object_id"`
	ObjectType string      `json:"object_type"`
	EntityType string      `json:"entity_type"`
	Entity     string      `json:"entity"`
	ReadID     []uuid.UUID `json:"read_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
