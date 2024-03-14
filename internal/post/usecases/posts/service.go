package posts

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/pkg/events"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	postRepo           PostRepo
	redis              redis.RedisEngine
	notiEventPublisher NotiEventPublisher
	authDomainService  domain.AuthDomainService
	groupDomainService domain.GroupDomainService
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	postRepo PostRepo,
	redis redis.RedisEngine,
	notiEventPublisher NotiEventPublisher,
	authDomainService domain.AuthDomainService,
	groupDomainService domain.GroupDomainService,
) UseCase {
	return &service{
		postRepo:           postRepo,
		redis:              redis,
		notiEventPublisher: notiEventPublisher,
		authDomainService:  authDomainService,
		groupDomainService: groupDomainService,
	}
}

// GetPostsByFeedGroup implements UseCase.
func (uc *service) GetPostsByFeedGroup(ctx context.Context, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	key := ""
	for _, groupId := range groupIds {
		key += groupId.String() + "_"
	}
	keyCache := constant.CachePostsFeedGroup + key + constant.CacheLimit +
		utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByGroupId(ctx, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
		}
		err = uc.redis.Set(keyCache, posts)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByFeed implements UseCase.
func (uc *service) GetPostsByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	key := ""
	for _, userId := range userIds {
		key += userId.String() + "_"
	}
	for _, groupId := range groupIds {
		key += groupId.String() + "_"
	}
	keyCache := constant.CachePostsFeed + key + constant.CacheLimit +
		utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByFeed(ctx, userIds, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
		}
		err = uc.redis.Set(keyCache, posts)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByGroupId implements UseCase.
func (uc *service) GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	groupIds := make([]uuid.UUID, 0)
	groupIds = append(groupIds, groupId)
	var posts []*domain.Post
	keyCache := constant.CachePostsGroupId + groupId.String() + constant.CacheLimit +
		utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByGroupId(ctx, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
		}
		err = uc.redis.Set(keyCache, posts)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByUserId implements UseCase.
func (uc *service) GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	keyCache := constant.CachePostsUserId + userId.String() + constant.CacheLimit +
		utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByUserId(ctx, userId, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByUserId failed")
		}
		err = uc.redis.Set(keyCache, posts)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// CreatePost implements UseCase.
func (uc *service) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Create(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create")
	}
	err = uc.redis.InvalidatePrefix(constant.CachePostsFeed)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	if post.GroupID.UUID.String() != constant.NullUUID {
		err = uc.redis.InvalidatePrefix(constant.CachePostsGroupId + post.GroupID.UUID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
		err = uc.redis.InvalidatePrefix(constant.CachePostsFeedGroup)
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	} else {
		err = uc.redis.InvalidatePrefix(constant.CachePostsUserId + post.UserID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	}
	eventPublish(ctx, uc, post)
	return post, nil
}
func eventPublish(ctx context.Context, uc *service, post *domain.Post) {
	var senderIds []string
	var typeNoti string
	data := map[string]interface{}{
		"content": post.Content,
		"userId":  post.UserID.String(),
		"groupId": post.GroupID.UUID.String(),
	}
	if post.GroupID.UUID.String() != constant.NullUUID {
		typeNoti = "group"
		res, err := uc.groupDomainService.GetGroupMembers(ctx, post.GroupID)
		if err != nil {
			errors.Wrap(err, "GetGroupMembers failed")
		}
		for _, member := range res.GroupMembers {
			senderIds = append(senderIds, member.UserId)
		}
		resGroup, err := uc.groupDomainService.GetGroup(ctx, post.GroupID)
		if err != nil {
			errors.Wrap(err, "GetGroup failed")
		}
		data["groupName"] = resGroup.Group.Name
	} else {
		userIds, err := uc.authDomainService.GetUserIdsByUserId(ctx, post.UserID)
		if err != nil {
			errors.Wrap(err, "GetUserIdsByUserId failed")
		}
		for _, userId := range userIds {
			senderIds = append(senderIds, userId.String())
		}
	}
	event := events.Noti{
		ActorID:    post.UserID.String(),
		SenderIDs:  utils.UniqueSlice(senderIds),
		Data:       data,
		Type:       typeNoti,
		ObjectType: "post",
		ObjectID:   post.ID.String(),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Marshal event failed")
	}
	uc.notiEventPublisher.Publish(ctx, eventBytes, "text/plain")
}

// DeletePost implements UseCase.
func (uc *service) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	isDelete, err := uc.postRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "postRepo.Delete")
	}
	err = uc.redis.InvalidatePrefix(constant.CachePosts)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	return isDelete, nil
}

// GetPost implements UseCase.
func (uc *service) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	post, err := uc.postRepo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Get")
	}
	return post, nil
}

// UpdatePost implements UseCase.
func (uc *service) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Update(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.UpdatePost")
	}
	err = uc.redis.InvalidatePrefix(constant.CachePostsFeed)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	if post.GroupID.UUID.String() != constant.NullUUID {
		err = uc.redis.InvalidatePrefix(constant.CachePostsGroupId + post.GroupID.UUID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
		err = uc.redis.InvalidatePrefix(constant.CachePostsFeedGroup)
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	} else {
		err = uc.redis.InvalidatePrefix(constant.CachePostsUserId + post.UserID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	}
	return post, nil
}
