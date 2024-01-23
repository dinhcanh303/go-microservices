package domain

import (
	"context"

	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
)

type UploadDomainService interface {
	GetAvatarUser(ctx context.Context, userId uuid.UUID) (*domainUpload.Attachment, error)
}
type GroupDomainService interface {
	GetGroupMembers(ctx context.Context, groupId string) (*gen.GetGroupMembersResponse, error)
}
