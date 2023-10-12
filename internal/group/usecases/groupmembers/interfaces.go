package groupmembers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/google/uuid"
)

type GroupMemberRepo interface {
	CreateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error)
	UpdateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id uuid.UUID) error
	DeleteAllGroupMembersByGroupId(ctx context.Context, groupId uuid.UUID) error
	GetAllGroupMembers(ctx context.Context, groupId uuid.UUID) ([]*domain.GroupMember, error)
	CountGroupMembers(ctx context.Context, groupId uuid.UUID) (int64, error)
}

type UseCase interface {
	CreateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error)
	UpdateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id uuid.UUID) error
	DeleteAllGroupMembersByGroupId(ctx context.Context, groupId uuid.UUID) error
	GetAllGroupMembers(ctx context.Context, groupId uuid.UUID) ([]*domain.GroupMember, error)
	CountGroupMembers(ctx context.Context, groupId uuid.UUID) (int64, error)
}
