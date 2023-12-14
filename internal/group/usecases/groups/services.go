package groups

import (
	"context"
	"encoding/json"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	groupmembers "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"golang.org/x/exp/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

type service struct {
	repo                 GroupRepo
	repoGroupMember      groupmembers.GroupMemberRepo
	groupCreatedEventPub GroupCreatedEventPublisher
	groupDeletedEventPub GroupDeletedEventPublisher
}

// GetAllGroupByUserId implements UseCase.
func (s *service) GetAllGroupByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Group, error) {
	result, err := s.repo.GetAllGroupByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllGroupByUserId")
	}
	return result, nil
}

// GetAllGroupByUserId implements UseCase.
func (s *service) GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]string, error) {
	result, err := s.repo.GetAllGroupIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetAllGroupIdByUserId")
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
		Role:    constant.OWNER,
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
	// err = s.repoGroupMember.DeleteAllGroupMembersByGroupId(ctx, id)
	// if err != nil {
	// 	slog.Error("service.DeleteGroup can't DeleteAllGroupMembers please check", err)
	// }
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

	result, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetGroup")
	}
	return result, nil
}

// Update implements UseCase.
func (s *service) UpdateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	result, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateGroup")
	}
	return result, nil
}

func NewService(repo GroupRepo, repoGroupMember groupmembers.GroupMemberRepo, groupCreatedEventPub GroupCreatedEventPublisher, groupDeletedEventPub GroupDeletedEventPublisher) UseCase {
	return &service{
		repo:                 repo,
		repoGroupMember:      repoGroupMember,
		groupCreatedEventPub: groupCreatedEventPub,
		groupDeletedEventPub: groupDeletedEventPub,
	}
}
