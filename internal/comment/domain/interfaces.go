package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type (
	LikeDomainService interface {
		GetLikesByCommentID(ctx context.Context, commentId uuid.UUID) ([]*domain.Like, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error)
	}
)
