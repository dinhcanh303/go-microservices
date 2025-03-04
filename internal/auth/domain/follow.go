package domain

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerId  uuid.UUID `json:"follower_id"`
	FollowingId uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewFollow(followerId uuid.UUID, followingId uuid.UUID) *Follow {
	return &Follow{
		FollowerId:  followerId,
		FollowingId: followingId,
	}
}
