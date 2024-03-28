package domain

import (
	"context"

	v1a "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1g "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error)
		CountCommentByPostID(ctx context.Context, postId uuid.UUID) (int64, error)
	}
	LikeDomainService interface {
		GetLikesByPostID(ctx context.Context, postId uuid.UUID) (*domain.LikesInfo, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error)
	}
	GroupDomainService interface {
		GetGroupIdsByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
		GetGroup(ctx context.Context, id uuid.NullUUID) (*v1g.GetGroupResponse, error)
		GetGroupMembers(ctx context.Context, groupId uuid.NullUUID) (*v1g.GetGroupMembersResponse, error)
	}
	AuthDomainService interface {
		GetUserIdsByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
		GetProfile(ctx context.Context, id uuid.UUID) (*v1a.GetProfileResponse, error)
	}
)
