package groups

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	groupmembers "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"golang.org/x/exp/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo            GroupRepo
	repoGroupMember groupmembers.GroupMemberRepo
	redis           redis.RedisEngine
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(
	repo GroupRepo,
	repoGroupMember groupmembers.GroupMemberRepo,
	redis redis.RedisEngine,
) UseCase {
	return &service{
		repo:            repo,
		repoGroupMember: repoGroupMember,
		redis:           redis,
	}
}

// GetGroups implements UseCase.
func (s *service) GetGroups(ctx context.Context) ([]*domain.Group, error) {
	var groups []*domain.Group
	keyCache := constant.CacheGroups
	err := utils.HandleHitCache(groups, s.redis, keyCache)
	if err != nil {
		groups, err = s.repo.GetGroups(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroups")
		}
		err = s.redis.Set(keyCache, groups)
		if err != nil {
			slog.Warn("failed set value in cache", err)
		}
	}
	return groups, nil
}

// GetGroupsByUserId implements UseCase.
func (s *service) GetGroupsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Group, error) {
	var groups []*domain.Group
	keyCache := constant.CacheGroupsByUserId + userId.String() +
		constant.CacheLimit + utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(groups, s.redis, keyCache)
	if err != nil {
		groups, err = s.repo.GetGroupsByUserId(ctx, userId, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroupsByUserId")
		}
		err = s.redis.Set(keyCache, groups)
		if err != nil {
			slog.Warn("failed set value in cache", err)
		}
	}

	return groups, nil
}

// GetGroupIdsByUserId implements UseCase.
func (s *service) GetGroupIdsByUserId(ctx context.Context, userId uuid.UUID) ([]string, error) {
	var groupIds []string
	keyCache := constant.CacheGroupIdsByUserId + userId.String()
	err := utils.HandleHitCache(groupIds, s.redis, keyCache)
	if err != nil {
		groupIds, err = s.repo.GetGroupIdsByUserId(ctx, userId)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroupIdsByUserId")
		}
		err = s.redis.Set(keyCache, groupIds)
		if err != nil {
			slog.Warn("failed set value in cache", err)
		}
	}
	return groupIds, nil
}

// Create implements UseCase.
func (s *service) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	result, err := s.repo.Create(ctx, &domain.Group{
		Name:        group.Name,
		Description: group.Description,
		Status:      group.Status,
		UserID:      user.ID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "service.Create")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroups)
	if err != nil {
		slog.Error("Invalidate cache group failed")
	}
	err = s.redis.Invalidate(constant.CacheGroup + group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
	}
	err = s.redis.Invalidate(constant.CacheGroupsByUserId + user.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group by user id failed : ID", user.ID.String())
	}
	err = s.redis.Invalidate(constant.CacheGroupIdsByUserId + user.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group ids by user id failed : ID", user.ID.String())
	}

	//Create the group member
	_, err = s.repoGroupMember.CreateGroupMember(ctx, &domain.GroupMember{
		ID:      uuid.New(),
		GroupID: result.ID,
		UserID:  user.ID,
		Role:    int32(sharedkernel.OWNER),
	})
	if err != nil {
		return nil, errors.New("create group member owner failed")
	}
	return result, nil
}

// Delete implements UseCase.
func (s *service) DeleteGroup(ctx context.Context, id uuid.UUID) (bool, error) {
	group, err := s.repo.Get(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteGroup")
	}
	result, err := s.repo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteGroup")
	}
	err = s.repoGroupMember.DeleteGroupMembersByGroupId(ctx, id)
	if err != nil {
		slog.Error("service.DeleteGroup can't DeleteAllGroupMembers please check", err)
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroups)
	if err != nil {
		slog.Error("Invalidate cache group failed")
	}
	err = s.redis.Invalidate(constant.CacheGroup + id.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", id)
	}
	err = s.redis.Invalidate(constant.CacheGroupsByUserId + group.UserID.String())
	if err != nil {
		slog.Error("Invalidate cache group by user id failed : ID", group.UserID.String())
	}
	err = s.redis.Invalidate(constant.CacheGroupIdsByUserId + group.UserID.String())
	if err != nil {
		slog.Error("Invalidate cache group ids by user id failed : ID", group.UserID.String())
	}
	return result, nil
}

// Get implements UseCase.
func (s *service) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	var group *domain.Group
	keyCache := constant.CacheGroup + id.String()
	err := utils.HandleHitCache(group, s.redis, keyCache)
	if err != nil {
		group, err = s.repo.Get(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroup")
		}
		err = s.redis.Set(keyCache, group)
		if err != nil {
			slog.Warn("failed set value in cache", err)
		}
	}
	return group, nil
}

// Update implements UseCase.
func (s *service) UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	group, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateGroup")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroups)
	if err != nil {
		slog.Error("Invalidate cache group failed")
	}
	err = s.redis.Invalidate(constant.CacheGroup + group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
	}
	err = s.redis.Invalidate(constant.CacheGroupsByUserId + group.UserID.String())
	if err != nil {
		slog.Error("Invalidate cache group by user id failed : ID", group.UserID.String())
	}
	err = s.redis.Invalidate(constant.CacheGroupIdsByUserId + group.UserID.String())
	if err != nil {
		slog.Error("Invalidate cache group ids by user id failed : ID", group.UserID.String())
	}
	return group, nil
}
