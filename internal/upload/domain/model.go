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
	URL            string    `json:"url"`
	URLThumbnail   string    `json:"url_thumbnail"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewAttachment(
	userId uuid.UUID,
	fileName,
	extension,
	mimeType,
	folder,
	urlThumbnail,
	url string) *Attachment {
	return &Attachment{
		ID:           uuid.New(),
		UserID:       userId,
		FileName:     fileName,
		Extension:    extension,
		MimeType:     mimeType,
		Folder:       folder,
		URL:          url,
		URLThumbnail: urlThumbnail,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
