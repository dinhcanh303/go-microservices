package groups

import (
	"context"
	"encoding/json"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	groupmembers "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"golang.org/x/exp/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo                 GroupRepo
	repoGroupMember      groupmembers.GroupMemberRepo
	redis                redis.RedisEngine
	groupCreatedEventPub GroupCreatedEventPublisher
	groupDeletedEventPub GroupDeletedEventPublisher
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(
	repo GroupRepo,
	repoGroupMember groupmembers.GroupMemberRepo,
	redis redis.RedisEngine,
	groupCreatedEventPub GroupCreatedEventPublisher,
	groupDeletedEventPub GroupDeletedEventPublisher) UseCase {
	return &service{
		repo:                 repo,
		repoGroupMember:      repoGroupMember,
		redis:                redis,
		groupCreatedEventPub: groupCreatedEventPub,
		groupDeletedEventPub: groupDeletedEventPub,
	}
}

var CACHE_SV_GROUP_GROUP_ID = "sv_group_group_id_"

// GetGroupsByUserId implements UseCase.
func (s *service) GetGroupsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Group, error) {
	result, err := s.repo.GetGroupsByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetGroupsByUserId")
	}
	return result, nil
}

// GetGroupIdsByUserId implements UseCase.
func (s *service) GetGroupIdsByUserId(ctx context.Context, userId uuid.UUID) ([]string, error) {
	result, err := s.repo.GetGroupIdsByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetGroupIdsByUserId")
	}
	return result, nil
}

// Create implements UseCase.
func (s *service) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	slog.Info("Create Group Service")
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
	err = s.redis.Invalidate(CACHE_SV_GROUP_GROUP_ID + group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
	}
	if result != nil {
		// Publish event created group
		eventBytes, err := json.Marshal(event.GroupCreated{
			ID:     result.ID,
			Name:   result.Name,
			Avatar: result.Description,
		})
		if err != nil {
			slog.Error("json marshal error", err)
		}
		s.groupCreatedEventPub.Publish(ctx, eventBytes, "text/plain")
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
	result, err := s.repo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteGroup")
	}
	err = s.repoGroupMember.DeleteGroupMembersByGroupId(ctx, id)
	if err != nil {
		slog.Error("service.DeleteGroup can't DeleteAllGroupMembers please check", err)
	}
	//Del cache
	err = s.redis.Invalidate(CACHE_SV_GROUP_GROUP_ID + id.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", id)
	}
	if result {
		// Publish event created group
		eventBytes, err := json.Marshal(event.GroupDeleted{
			ID: id,
		})
		if err != nil {
			slog.Error("json marshal error", err)
		}
		s.groupDeletedEventPub.Publish(ctx, eventBytes, "text/plain")
	}
	return result, nil
}

// Get implements UseCase.
func (s *service) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	var group *domain.Group
	keyCache := CACHE_SV_GROUP_GROUP_ID + id.String()
	err := utils.HandleHitCache(group, s.redis, keyCache)
	if err != nil {
		slog.Info("MISS_CACHE", err)
		group, err = s.repo.Get(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroup")
		}
		err = s.redis.Set(keyCache, group, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}

	return group, nil
}

// Update implements UseCase.
func (s *service) UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	result, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateGroup")
	}
	//Del cache
	err = s.redis.Invalidate(CACHE_SV_GROUP_GROUP_ID + group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
	}
	return result, nil
}
