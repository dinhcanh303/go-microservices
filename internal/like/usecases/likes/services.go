package likes

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	likeRepo LikeRepo
	redis    redis.RedisEngine
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(
	likeRepo LikeRepo,
	redis redis.RedisEngine) UseCase {
	return &service{
		likeRepo: likeRepo,
		redis:    redis,
	}
}

// GetLikeByUserId implements UseCase.
func (s *service) GetLikeByUserId(ctx context.Context, likeableType string, likeableId uuid.UUID, userId uuid.UUID) (*domain.Like, error) {
	like, err := s.likeRepo.GetLikeByUserId(ctx, likeableType, likeableId, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllLikeByCommentID")
	}
	return like, nil
}

// GetLikesInfoByCommentID implements UseCase.
func (s *service) GetLikesInfoByCommentID(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (*domain.LikesInfo, error) {
	var likeInfo *domain.LikesInfo
	keyCache := constant.CacheLikeInfoByLikeableId + commentID.String() + constant.CacheUserId + userID.String()
	err := utils.HandleHitCache(likeInfo, s.redis, keyCache)
	if err != nil {
		likeInfo, err = s.likeRepo.GetLikesInfoByCommentID(ctx, commentID, userID)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetLikesInfoByCommentID")
		}
		err = s.redis.Set(keyCache, likeInfo)
		if err != nil {
			slog.Error("set cache like info by comment id", err)
		}
	}
	return likeInfo, nil
}

// GetLikesInfoByPostID implements UseCase.
func (s *service) GetLikesInfoByPostID(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (*domain.LikesInfo, error) {
	var likeInfo *domain.LikesInfo
	keyCache := constant.CacheLikeInfoByLikeableId + postID.String() + constant.CacheUserId + userID.String()
	err := utils.HandleHitCache(likeInfo, s.redis, keyCache)
	if err != nil {
		likeInfo, err = s.likeRepo.GetLikesInfoByPostID(ctx, postID, userID)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetLikesInfoByPostID")
		}
		err = s.redis.Set(keyCache, likeInfo)
		if err != nil {
			slog.Error("set cache like info by comment id", err)
		}
	}

	return likeInfo, nil
}

// CreateLike implements UseCase.
func (s *service) CreateLike(ctx context.Context, like *domain.Like) (*domain.Like, error) {
	like, err := s.likeRepo.Create(ctx, like)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateLike")
	}
	//Del cache
	keyCache := constant.CacheLikeInfoByLikeableId + like.LikeableID.String() +
		constant.CacheUserId + like.UserID.String()
	err = s.redis.Invalidate(keyCache)
	if err != nil {
		slog.Error("Invalidate cache like info failed", err)
	}
	return like, nil
}

// DeleteLike implements UseCase.
func (s *service) DeleteLike(ctx context.Context, id uuid.UUID) (bool, error) {
	result, err := s.likeRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.CreateLike")
	}
	//Del cache
	keyCache := constant.CacheLikeInfoByLikeableId
	err = s.redis.InvalidatePrefix(keyCache)
	if err != nil {
		slog.Error("Invalidate cache like info failed", err)
	}
	return result, nil
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
	//Del cache
	keyCache := constant.CacheLikeInfoByLikeableId + like.LikeableID.String() +
		constant.CacheUserId + like.UserID.String()
	err = s.redis.Invalidate(keyCache)
	if err != nil {
		slog.Error("Invalidate cache like info failed", err)
	}
	return like, nil
}
