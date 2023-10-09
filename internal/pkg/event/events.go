package event

import (
	shared "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type GroupUploadBegin struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
}

func (e GroupUploadBegin) Identity() string {
	return "GroupUploadBegin"
}

type GroupUploadEnd struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
	Url            string    `json:"url"`
	UrlThumbnail   string    `json:"url_thumbnail"`
}

func (e *GroupUploadEnd) Identity() string {
	return "GroupUploadEnd"
}

type PostUploadBegin struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
}

func (e PostUploadBegin) Identity() string {
	return "PostUploadBegin"
}

type PostUploadEnd struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
	Url            string    `json:"url"`
	UrlThumbnail   string    `json:"url_thumbnail"`
}

func (e *PostUploadEnd) Identity() string {
	return "PostUploadEnd"
}

type CommentUploadBegin struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
}

func (e CommentUploadBegin) Identity() string {
	return "CommentUploadBegin"
}

type CommentUploadEnd struct {
	shared.DomainEvent
	UserID         uuid.UUID `json:"user_id"`
	AttachableType string    `json:"attachable_type"`
	AttachableID   uuid.UUID `json:"attachable_id"`
	Url            string    `json:"url"`
	UrlThumbnail   string    `json:"url_thumbnail"`
}

func (e *CommentUploadEnd) Identity() string {
	return "CommentUploadEnd"
}
