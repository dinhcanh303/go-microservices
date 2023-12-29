package domain

import (
	"context"

	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type UploadDomainService interface {
	GetAvatarUser(ctx context.Context, userId uuid.UUID) (*domainUpload.Attachment, error)
}
