package domain

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID             uuid.UUID `json:"id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
	UserID         uuid.UUID `json:"user_id"`
	FileName       string    `json:"filename"`
	Extension      string    `json:"extension"`
	MimeType       string    `json:"mime_type"`
	Folder         string    `json:"folder"`
	URLThumbnail   string    `json:"url_thumbnail"`
	URLType        string    `json:"url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
