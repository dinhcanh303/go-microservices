package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type orderRepo struct {
	pg postgres.DBEngine
}

func NewGroupRepo(pg postgres.DBEngine) groups.GroupRepo {
	return &orderRepo{pg: pg}
}

var _ groups.GroupRepo = (*orderRepo)(nil)

var RepositorySet = wire.NewSet(NewGroupRepo)

// Create implements groups.GroupRepo.
func (rp *orderRepo) Create(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "CreateGroupRepo")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:          group.ID,
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
func (rp *orderRepo) Delete(ctx context.Context, uuid string) (bool, error) {
	panic("unimplemented")
}

// Get implements groups.GroupRepo.
func (rp *orderRepo) Get(ctx context.Context, uuid string) (*domain.Group, error) {
	panic("unimplemented")
}

// GetWithUnscoped implements groups.GroupRepo.
func (rp *orderRepo) GetWithUnscoped(ctx context.Context, uuid string) (*domain.Group, error) {
	panic("unimplemented")
}

// Update implements groups.GroupRepo.
func (rp *orderRepo) Update(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	panic("unimplemented")
}
