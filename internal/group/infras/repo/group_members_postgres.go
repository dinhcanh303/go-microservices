package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/infras/postgresql"
	groupmembers "github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type groupMemberRepo struct {
	pg postgres.DBEngine
}

// CountGroupMember implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) CountGroupMembers(ctx context.Context, groupId uuid.UUID) (int64, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "CreateGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.CountGroupMembers(ctx, groupId)
	if err != nil {
		return 0, errors.Wrap(err, "qtx.Create(ctx, groupId) failed")
	}
	return result, tx.Commit()
}

// CreateGroupMember implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) CreateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreateGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.CreateGroupMember(ctx, postgresql.CreateGroupMemberParams{
		ID:      uuid.New(),
		GroupID: groupMember.GroupID,
		UserID:  groupMember.UserID,
		Role:    groupMember.Role,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateGroupMemberParams) failed")
	}

	return &domain.GroupMember{
		ID:        result.ID,
		GroupID:   result.GroupID,
		UserID:    result.UserID,
		Role:      result.Role,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, tx.Commit()
}

// DeleteGroupMembersByGroupId implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) DeleteGroupMembersByGroupId(ctx context.Context, groupId uuid.UUID) error {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Repo DeleteGroupMembersByGroupId db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteGroupMembersByGroupId(ctx, groupId)
	if err != nil {
		return errors.Wrap(err, "qtx.DeleteGroupMembersByGroupId(ctx, groupId) failed")
	}
	return tx.Commit()
}

// DeleteGroupMember implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) DeleteGroupMember(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "Repo DeleteGroupMember db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteGroupMember(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "qtx.DeleteGroupMember(ctx, ids) failed")
	}
	return true, tx.Commit()
}

// GetGroupMembers implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) GetGroupMembers(ctx context.Context, groupId uuid.UUID) ([]*domain.GroupMember, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	results, err := querier.GetGroupMembers(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.GetGroupMembers(ctx, groupId) failed")
	}
	return lo.Map(results, func(item postgresql.GroupGroupMember, _ int) *domain.GroupMember {
		return &domain.GroupMember{
			ID:        item.ID,
			GroupID:   item.GroupID,
			UserID:    item.UserID,
			Role:      item.Role,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}), nil
}

// UpdateGroupMember implements groupmembers.GroupMemberRepo.
func (rp *groupMemberRepo) UpdateGroupMember(ctx context.Context, groupMember *domain.GroupMember) (*domain.GroupMember, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "Repo UpdateGroupMember db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.UpdateGroupMember(ctx, postgresql.UpdateGroupMemberParams{
		ID: groupMember.ID,
		Role: sql.NullInt32{
			Int32: groupMember.Role,
			Valid: groupMember.Role != 0,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.UpdateGroupMember(ctx, postgresql.UpdateGroupMemberParams) failed")
	}
	return &domain.GroupMember{
		ID:        result.ID,
		GroupID:   result.GroupID,
		UserID:    result.UserID,
		Role:      result.Role,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, tx.Commit()
}

func NewGroupMemberRepo(pg postgres.DBEngine) groupmembers.GroupMemberRepo {
	return &groupMemberRepo{pg: pg}
}

var _ groupmembers.GroupMemberRepo = (*groupMemberRepo)(nil)

var RepositoryGroupMemberSet = wire.NewSet(NewGroupMemberRepo)
