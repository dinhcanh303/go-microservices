package event

import shared "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"

type AttachmentUpload struct {
	shared.DomainEvent
}

func (e AttachmentUpload) Identity() string {
	return "AttachmentUpload"
}
