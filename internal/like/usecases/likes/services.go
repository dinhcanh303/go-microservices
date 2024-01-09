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

// GetLikesInfoByCommentID implements UseCase.
func (s *service) GetLikesInfoByCommentID(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (*domain.LikesInfo, error) {
	likeInfo, err := s.likeRepo.GetLikesInfoByCommentID(ctx, commentID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetLikesInfoByCommentID")
	}
	return likeInfo, nil
}

// GetLikesInfoByPostID implements UseCase.
func (s *service) GetLikesInfoByPostID(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (*domain.LikesInfo, error) {
	likeInfo, err := s.likeRepo.GetLikesInfoByPostID(ctx, postID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetLikesInfoByPostID")
	}
	return likeInfo, nil
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
func (s *service) GetLikesByCommentID(ctx context.Context, commentID uuid.UUID) ([]*domain.Like, error) {
	likes, err := s.likeRepo.GetLikesByCommentID(ctx, commentID)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllLikeByCommentID")
	}
	return likes, nil
}

// GetAllLikeByPostID implements UseCase.
func (s *service) GetLikesByPostID(ctx context.Context, postID uuid.UUID) ([]*domain.Like, error) {
	likes, err := s.likeRepo.GetLikesByPostID(ctx, postID)
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
