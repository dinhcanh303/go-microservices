package likes

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	likeRepo LikeRepo
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(likeRepo LikeRepo) UseCase {
	return &service{
		likeRepo: likeRepo,
	}
}

// CreateLike implements UseCase.
func (s *service) CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	like, err := s.likeRepo.Create(ctx, like)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateLike")
	}
	return like, nil
}

// DeleteLike implements UseCase.
func (s *service) DeleteLike(ctx context.Context, id uuid.UUID) (bool, error) {
	like, err := s.likeRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.CreateLike")
	}
	return like, nil
}

// GetAllLikeByCommentID implements UseCase.
func (s *service) GetAllLikeByCommentID(ctx context.Context, commentID uuid.UUID) ([]*domain.Like, error) {
	likes, err := s.likeRepo.GetAllByCommentID(ctx, commentID)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllLikeByCommentID")
	}
	return likes, nil
}

// GetAllLikeByPostID implements UseCase.
func (s *service) GetAllLikeByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Like, error) {
	likes, err := s.likeRepo.GetAllByPostID(ctx, postID)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllLikeByPostID")
	}
	return likes, nil
}

// UpdateLike implements UseCase.
func (s *service) UpdateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	like, err := s.likeRepo.Update(ctx, like)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateLike")
	}
	return like, nil
}
