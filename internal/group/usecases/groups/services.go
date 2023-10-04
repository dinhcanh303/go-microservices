package groups

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"

	"github.com/google/wire"
	"github.com/pkg/errors"
)

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

type service struct {
	repo domain.GroupRepo
}

// Create implements UseCase.
func (s *service) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	result, err := s.repo.Create(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.Create")
	}
	return result, nil
}

// Delete implements UseCase.
func (s *service) DeleteGroup(ctx context.Context, uuid string) (bool, error) {
	result, err := s.repo.Delete(ctx, uuid)
	if err != nil {
		return false, errors.Wrap(err, "service.Delete")
	}
	return result, nil
}

// Get implements UseCase.
func (s *service) GetGroup(ctx context.Context, uuid string) (*domain.Group, error) {
	result, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, errors.Wrap(err, "service.Get")
	}
	return result, nil
}

// GetWithUnscoped implements UseCase.
func (s *service) GetGroupWithUnscoped(ctx context.Context, uuid string) (*domain.Group, error) {
	result, err := s.repo.GetWithUnscoped(ctx, uuid)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetWithUnscoped")
	}
	return result, nil
}

// Update implements UseCase.
func (s *service) UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	result, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetWithUnscoped")
	}
	return result, nil
}

func NewService(repo domain.GroupRepo) UseCase {
	return &service{
		repo: repo,
	}
}
