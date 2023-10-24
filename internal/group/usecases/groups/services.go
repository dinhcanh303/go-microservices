package groups

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	groupmembers "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"golang.org/x/exp/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

type service struct {
	repo            GroupRepo
	repoGroupMember groupmembers.GroupMemberRepo
}

// GetAllGroupByUserId implements UseCase.
func (s *service) GetAllGroupByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Group, error) {
	result, err := s.repo.GetAllGroupByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllGroupByUserId")
	}
	return result, nil
}

// GetAllGroupByUserId implements UseCase.
func (s *service) GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]string, error) {
	result, err := s.repo.GetAllGroupIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllGroupIdByUserId")
	}
	return result, nil
}

// Create implements UseCase.
func (s *service) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	slog.Info("Create Group Service")
	result, err := s.repo.Create(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.Create")
	}
	return result, nil
}

// Delete implements UseCase.
func (s *service) DeleteGroup(ctx context.Context, id uuid.UUID) (bool, error) {
	result, err := s.repo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.Delete")
	}
	err = s.repoGroupMember.DeleteAllGroupMembersByGroupId(ctx, id)
	if err != nil {
		slog.Error("service.DeleteGroup can't DeleteAllGroupMembers please check", err)
	}
	return result, nil
}

// Get implements UseCase.
func (s *service) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	result, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "service.Get")
	}
	return result, nil
}

// Update implements UseCase.
func (s *service) UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	result, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateGroup")
	}
	return result, nil
}

func NewService(repo GroupRepo, repoGroupMember groupmembers.GroupMemberRepo) UseCase {
	return &service{
		repo:            repo,
		repoGroupMember: repoGroupMember,
	}
}
