package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type userRepo struct {
	pg postgres.DBEngine
}

// GetAllUserIdOfCompany implements auth.UserRepo.
func (rp *userRepo) GetAllUserIdOfCompany(ctx context.Context, company string) ([]uuid.UUID, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	userIds, err := querier.GetAllUserIdOfCompany(ctx, sql.NullString{
		String: company,
		Valid:  company != "",
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetAllUserIdOfCompany failed")
	}
	return userIds, nil
}

// CreateUser implements auth.AuthRepo.
func (rp *userRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.CreateUser(ctx, postgresql.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName: sql.NullString{
			String: user.FullName,
			Valid:  user.FullName != "",
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.Create failed")
	}
	return &domain.User{
		ID:        result.ID,
		Email:     result.Email,
		Password:  result.Password,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, tx.Commit()
}

// GetUser implements auth.AuthRepo.
func (rp *userRepo) GetUser(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	user, err := querier.GetUser(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUser failed")
	}
	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FullName.String,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserByEmail implements auth.AuthRepo.
func (rp *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	user, err := querier.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUserByEmail failed")
	}
	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FullName.String,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

var _ auth.UserRepo = (*userRepo)(nil)

func NewUserRepo(pg postgres.DBEngine) auth.UserRepo {
	return &userRepo{pg: pg}
}

var UserRepoSet = wire.NewSet(NewUserRepo)
