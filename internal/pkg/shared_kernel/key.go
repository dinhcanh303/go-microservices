package sharedkernel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Key struct {
	ID                int64           `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	PublicKey         string          `json:"public_key"`
	PrivateKey        string          `json:"private_key"`
	RefreshToken      string          `json:"refresh_token"`
	RefreshTokensUsed json.RawMessage `json:"refresh_tokens_used"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
