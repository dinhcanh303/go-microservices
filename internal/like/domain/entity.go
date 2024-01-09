package domain

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID           uuid.UUID `json:"id"`
	Emoji        string    `json:"emoji"`
	LikeableType string    `json:"likeable_type"`
	LikeableID   uuid.UUID `json:"likeable_id"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type LikesInfo struct {
	YourLikedEmoji    string   `json:"your_liked_emoji"`
	YourLike          int64    `json:"your_like"`
	OthersLikedEmojis []string `json:"others_liked_emojis"`
	OthersLikes       int64    `json:"others_likes"`
}
