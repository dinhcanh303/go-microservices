package posts

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type usecase struct {
	postRepo PostRepo
	redis    redis.RedisEngine
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	postRepo PostRepo,
	redis redis.RedisEngine,
) UseCase {
	return &usecase{
		postRepo: postRepo,
		redis:    redis,
	}
}

var (
	CACHE_SV_POST            = "sv_post_"
	CACHE_SV_POST_FEED       = "sv_post_feed_"
	CACHE_SV_POST_FEED_GROUP = "sv_post_feed_group_"
	CACHE_SV_POST_GROUP_ID   = "sv_post_group_id_"
	CACHE_SV_POST_USER_ID    = "sv_post_user_id_"
)

// GetPostsByFeedGroup implements UseCase.
func (uc *usecase) GetPostsByFeedGroup(ctx context.Context, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	key := ""
	for _, groupId := range groupIds {
		key += groupId.String() + "_"
	}
	keyCache := CACHE_SV_POST_FEED_GROUP + key + utils.String(limit) + "_" + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByGroupId(ctx, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
		}
		err = uc.redis.Set(keyCache, posts, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByFeed implements UseCase.
func (uc *usecase) GetPostsByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	key := ""
	for _, userId := range userIds {
		key += userId.String() + "_"
	}
	for _, groupId := range groupIds {
		key += groupId.String() + "_"
	}
	keyCache := CACHE_SV_POST_FEED + key + utils.String(limit) + "_" + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByFeed(ctx, userIds, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
		}
		err = uc.redis.Set(keyCache, posts, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByGroupId implements UseCase.
func (uc *usecase) GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	groupIds := make([]uuid.UUID, 0)
	groupIds = append(groupIds, groupId)
	var posts []*domain.Post
	keyCache := CACHE_SV_POST_GROUP_ID + groupId.String() + "_" + utils.String(limit) + "_" + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByGroupId(ctx, groupIds, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
		}
		err = uc.redis.Set(keyCache, posts, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// GetPostsByUserId implements UseCase.
func (uc *usecase) GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	var posts []*domain.Post
	keyCache := CACHE_SV_POST_USER_ID + userId.String() + "_" + utils.String(limit) + "_" + utils.String(offset)
	err := utils.HandleHitCache(posts, uc.redis, keyCache)
	if err != nil {
		posts, err = uc.postRepo.GetByUserId(ctx, userId, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetPostsByUserId failed")
		}
		err = uc.redis.Set(keyCache, posts, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return posts, nil
}

// CreatePost implements UseCase.
func (uc *usecase) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Create(ctx, post)
	slog.Info("CREATE POST", post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create")
	}
	err = uc.redis.InvalidatePrefix(CACHE_SV_POST_FEED)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	if post.GroupID.UUID.String() != constant.NullUUID {
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_GROUP_ID + post.GroupID.UUID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_FEED_GROUP)
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	} else {
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_USER_ID + post.UserID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	}
	return post, nil
}

// DeletePost implements UseCase.
func (uc *usecase) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	isDelete, err := uc.postRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "postRepo.Delete")
	}
	err = uc.redis.InvalidatePrefix(CACHE_SV_POST)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	return isDelete, nil
}

// GetPost implements UseCase.
func (uc *usecase) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	post, err := uc.postRepo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Get")
	}
	return post, nil
}

// UpdatePost implements UseCase.
func (uc *usecase) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Update(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.UpdatePost")
	}
	err = uc.redis.InvalidatePrefix(CACHE_SV_POST_FEED)
	if err != nil {
		slog.Error("InvalidatePrefix cache key failed")
	}
	if post.GroupID.UUID.String() != constant.NullUUID {
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_GROUP_ID + post.GroupID.UUID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_FEED_GROUP)
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	} else {
		err = uc.redis.InvalidatePrefix(CACHE_SV_POST_USER_ID + post.UserID.String())
		if err != nil {
			slog.Error("InvalidatePrefix cache key failed")
		}
	}
	return post, nil
}
