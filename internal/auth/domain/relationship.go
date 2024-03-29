package domain

import (
	"time"

	"github.com/google/uuid"
)

type Relationship struct {
	ID         int64     `json:"id"`
	FollowerId uuid.UUID `json:"follower_id"`
	FollowedId uuid.UUID `json:"followed_id"`
	Status     int32     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
