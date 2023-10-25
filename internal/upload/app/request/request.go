package request

import "github.com/google/uuid"

type GetAttachmentRequest struct {
	AttachmentID uuid.UUID `json:"id"`
}
type UpdateAttachmentRequest struct {
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
}
type DeleteAttachmentRequest struct {
	AttachmentID uuid.UUID `json:"id"`
}
