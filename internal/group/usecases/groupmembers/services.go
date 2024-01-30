package groupmembers

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo  GroupMemberRepo
	redis redis.RedisEngine
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(repo GroupMemberRepo, redis redis.RedisEngine) UseCase {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

// CheckGroupMember implements UseCase.
func (s *service) GetRoleOfGroupMember(ctx context.Context, groupId uuid.UUID, userId uuid.UUID) (int32, error) {
	result, err := s.repo.GetRoleOfGroupMember(ctx, groupId, userId)
	if err != nil {
		return 0, errors.Wrap(err, "service.CheckGroupMember")
	}
	return result, nil
}

// CreateGroupMember implements UseCase.
func (s *service) CreateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	result, err := s.repo.CreateGroupMember(ctx, groupMember)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateGroupMember")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroup + result.GroupID.String())
	if err != nil {
		slog.Warn("Invalidate cache group members failed", err)
	}
	return result, nil
}

// DeleteGroupMember implements UseCase.
func (s *service) DeleteGroupMember(ctx context.Context, id uuid.UUID) (bool, error) {
	deleted, err := s.repo.DeleteGroupMember(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteGroupMember")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroup + id.String())
	if err != nil {
		slog.Warn("Invalidate cache group members failed", err)
	}
	return deleted, nil
}

// UpdateGroupMember implements UseCase.
func (s *service) UpdateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	result, err := s.repo.UpdateGroupMember(ctx, groupMember)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateGroupMember")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroup + result.GroupID.String())
	if err != nil {
		slog.Warn("Invalidate cache group members failed", err)
	}
	return result, nil
}

// DeleteAllGroupMembersByGroupId implements UseCase.
func (s *service) DeleteGroupMembersByGroupId(ctx context.Context, groupId uuid.UUID) error {
	err := s.repo.DeleteGroupMembersByGroupId(ctx, groupId)
	if err != nil {
		return errors.Wrap(err, "service.DeleteGroupMembersByGroupId")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheGroup + groupId.String())
	if err != nil {
		slog.Warn("Invalidate cache group members failed", err)
	}
	return nil
}

// CountGroupMember implements UseCase.
func (s *service) CountGroupMembers(ctx context.Context, groupId uuid.UUID) (int64, error) {
	result, err := s.repo.CountGroupMembers(ctx, groupId)
	if err != nil {
		return 0, errors.Wrap(err, "service.CountGroupMember")
	}
	return result, nil
}

// GetAllGroupMember implements UseCase.
func (s *service) GetGroupMembers(ctx context.Context, groupId uuid.UUID) ([]*domain.GroupMember, error) {
	var groups []*domain.GroupMember
	keyCache := constant.CacheGroupMembers + groupId.String()
	err := utils.HandleHitCache(groups, s.redis, keyCache)
	if err != nil {
		groups, err = s.repo.GetGroupMembers(ctx, groupId)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetGroupMembers")
		}
		err = s.redis.Set(keyCache, groups)
		if err != nil {
			slog.Warn("failed to set key cache")
		}
	}
	return groups, nil
}
