package domain

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type UploadDomainService interface {
	GetAvatarUser(ctx context.Context, userId uuid.UUID) (*domainUpload.Attachment, error)
}
type GroupDomainService interface {
	GetGroupMembers(ctx context.Context, groupId string) (*v1.GetGroupMembersResponse, error)
}
