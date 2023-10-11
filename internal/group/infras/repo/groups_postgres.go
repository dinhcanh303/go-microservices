package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type groupRepo struct {
	pg postgres.DBEngine
}

func NewGroupRepo(pg postgres.DBEngine) groups.GroupRepo {
	return &groupRepo{pg: pg}
}

var _ groups.GroupRepo = (*groupRepo)(nil)

var RepositoryGroupSet = wire.NewSet(NewGroupRepo)

// Create implements groups.GroupRepo.
func (rp *groupRepo) Create(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	slog.Info("Repo Postgresql")
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreateGroupRepo")
	}
	qtx := querier.WithTx(tx)
	slog.Info("QTX")
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:          uuid.New(),
		Name:        group.Name,
		Description: group.Description,
		Status:      group.Status,
		UserID:      group.UserID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return &domain.Group{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Status:      result.Status,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		DeletedAt:   result.DeletedAt.Time,
	}, tx.Commit()
}

// Delete implements groups.GroupRepo.
func (rp *groupRepo) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "DeleteGroupRepo")
	}
	qtx := querier.WithTx(tx)
	err = qtx.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "qtx.Delete(ctx, id) failed")
	}
	return true, tx.Commit()
}

// Get implements groups.GroupRepo.
func (rp *groupRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "GetGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, id) failed")
	}

	return &domain.Group{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Status:      result.Status,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		DeletedAt:   result.DeletedAt.Time,
	}, tx.Commit()
}

// GetWithUnscoped implements groups.GroupRepo.
func (rp *groupRepo) GetWithUnscoped(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "GetGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.GetWithUnscoped(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.GetWithUnscoped(ctx, id) failed")
	}
	return &domain.Group{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Status:      result.Status,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		DeletedAt:   result.DeletedAt.Time,
	}, tx.Commit()
}

// Update implements groups.GroupRepo.
func (rp *groupRepo) Update(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreateGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Update(ctx, postgresql.UpdateParams{
		Name:        group.Name,
		Description: group.Description,
		Status:      group.Status,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return &domain.Group{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Status:      result.Status,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		DeletedAt:   result.DeletedAt.Time,
	}, tx.Commit()
}
