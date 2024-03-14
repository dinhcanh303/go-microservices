package likes

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/pkg/events"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	likeRepo           LikeRepo
	redis              redis.RedisEngine
	notiEventPublisher NotiEventPublisher
	postDomainSvc      domain.PostDomainService
	commentDomainSvc   domain.CommentDomainService
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(
	likeRepo LikeRepo,
	redis redis.RedisEngine,
	notiEventPublisher NotiEventPublisher,
	postDomainSvc domain.PostDomainService,
	commentDomainSvc domain.CommentDomainService,
) UseCase {
	return &service{
		likeRepo:           likeRepo,
		redis:              redis,
		notiEventPublisher: notiEventPublisher,
		postDomainSvc:      postDomainSvc,
		commentDomainSvc:   commentDomainSvc,
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
	eventPublish(ctx, s, like)
	return like, nil
}
func eventPublish(ctx context.Context, uc *service, like *domain.Like) {
	var senderIds []string
	var typeNoti string
	data := map[string]interface{}{
		"content":      like.Emoji,
		"userId":       like.UserID.String(),
		"likeableType": like.LikeableType,
		"likeableId":   like.LikeableID.String(),
	}
	if like.LikeableType == constant.LikePostType {
		res, _ := uc.postDomainSvc.GetPostNormal(ctx, like.LikeableID)
		if res.Post.GroupId != constant.NullUUID {
			typeNoti = "group"
		}
		data["postId"] = res.Post.Id
		data["groupId"] = res.Post.GroupId
		data["postContent"] = res.Post.Content
		senderIds = append(senderIds, res.Post.UserId)
	} else {
		res, _ := uc.commentDomainSvc.GetComment(ctx, like.LikeableID)
		data["postId"] = res.Comment.PostId
		data["parentCommentId"] = res.Comment.ParentCommentId
		data["commentContent"] = res.Comment.Content
		senderIds = append(senderIds, res.Comment.UserId)
	}
	event := events.Noti{
		ActorID:    like.UserID.String(),
		SenderIDs:  utils.UniqueSlice(senderIds),
		Data:       data,
		Type:       typeNoti,
		ObjectType: "like",
		ObjectID:   like.ID.String(),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Marshal event failed")
	}
	uc.notiEventPublisher.Publish(ctx, eventBytes, "text/plain")
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
