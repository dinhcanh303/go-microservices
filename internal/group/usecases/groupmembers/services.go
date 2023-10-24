package groupmembers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo GroupMemberRepo
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(repo GroupMemberRepo) UseCase {
	return &service{
		repo: repo,
	}
}

// CreateGroupMember implements UseCase.
func (s *service) CreateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	result, err := s.repo.CreateGroupMember(ctx, groupMember)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateGroupMember")
	}
	return result, nil
}

// DeleteGroupMember implements UseCase.
func (s *service) DeleteGroupMember(ctx context.Context, id uuid.UUID) (bool, error) {
	deleted, err := s.repo.DeleteGroupMember(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteGroupMember")
	}
	return deleted, nil
}

// UpdateGroupMember implements UseCase.
func (s *service) UpdateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	result, err := s.repo.UpdateGroupMember(ctx, groupMember)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateGroupMember")
	}
	return result, nil
}

// DeleteAllGroupMembersByGroupId implements UseCase.
func (s *service) DeleteAllGroupMembersByGroupId(ctx context.Context, groupId uuid.UUID) error {
	err := s.repo.DeleteAllGroupMembersByGroupId(ctx, groupId)
	if err != nil {
		return errors.Wrap(err, "service.DeleteAllGroupMembersByGroupId")
	}
	return nil
}

// CountGroupMember implements UseCase.
func (s *service) CountGroupMembers(ctx context.Context, groupId uuid.UUID) (int64, error) {
	result, err := s.repo.CountGroupMembers(ctx, groupId)
	if err != nil {
		return 0, errors.Wrap(err, "service.CountGroupMember")
	}
	return result, nil
}

// GetAllGroupMember implements UseCase.
func (s *service) GetAllGroupMembers(ctx context.Context, groupId uuid.UUID) ([]*domain.GroupMember, error) {
	result, err := s.repo.GetAllGroupMembers(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllGroupMember")
	}
	return result, nil
}
