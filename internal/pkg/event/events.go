package event

import shared "go-microservices/internal/pkg/shared_kernel"

type AttachmentUpload struct {
	shared.DomainEvent
}

func (e AttachmentUpload) Identity() string {
	return "AttachmentUpload"
}
