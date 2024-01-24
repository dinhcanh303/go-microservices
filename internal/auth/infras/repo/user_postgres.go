package repo

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type userRepo struct {
	pg postgres.DBEngine
}

// GetUsersBirthDayByCurrentDay implements auth.UserRepo.
func (rp *userRepo) GetUsersBirthDayByCurrentDay(ctx context.Context) ([]*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	users, err := querier.GetUsersBirthDayByCurrentDay(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUsersBirthDayByCurrentDay failed")
	}
	return lo.Map(users, func(user postgresql.AuthUser, _ int) *domain.User {
		return &domain.User{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName.String,
			NickName:    user.NickName.String,
			Role:        user.Role.String,
			AvatarUrl:   user.AvatarUrl.String,
			ProfileUrl:  user.ProfileUrl.String,
			Gender:      user.Gender.Bool,
			Phone:       user.Phone.String,
			Address:     user.Address.String,
			DateOfBirth: user.DateOfBirth.Time,
			Position:    user.Position.String,
			// Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}), nil
}

// GetUsersBirthDayByCurrentMonth implements auth.UserRepo.
func (rp *userRepo) GetUsersBirthDayByCurrentMonth(ctx context.Context) ([]*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	users, err := querier.GetUsersBirthDayByCurrentMonth(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUsersBirthDayByCurrentMonth failed")
	}
	return lo.Map(users, func(user postgresql.AuthUser, _ int) *domain.User {
		return &domain.User{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName.String,
			NickName:    user.NickName.String,
			Role:        user.Role.String,
			AvatarUrl:   user.AvatarUrl.String,
			ProfileUrl:  user.ProfileUrl.String,
			Gender:      user.Gender.Bool,
			Phone:       user.Phone.String,
			Address:     user.Address.String,
			DateOfBirth: user.DateOfBirth.Time,
			Position:    user.Position.String,
			// Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}), nil
}

// GetUsersInviteGroup implements auth.UserRepo.
func (rp *userRepo) GetUsersInviteGroup(ctx context.Context, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	users, err := querier.GetUsersInviteGroup(ctx, postgresql.GetUsersInviteGroupParams{
		Column1: groupIds,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUsers failed")
	}
	return lo.Map(users, func(user postgresql.AuthUser, _ int) *domain.User {
		return &domain.User{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName.String,
			NickName:    user.NickName.String,
			Role:        user.Role.String,
			AvatarUrl:   user.AvatarUrl.String,
			ProfileUrl:  user.ProfileUrl.String,
			Gender:      user.Gender.Bool,
			Phone:       user.Phone.String,
			Address:     user.Address.String,
			DateOfBirth: user.DateOfBirth.Time,
			Position:    user.Position.String,
			// Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}), nil
}

// UpdateUser implements auth.UserRepo.
func (rp *userRepo) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.UpdateUser db failed")
	}
	qtx := querier.WithTx(tx)
	slog.Info("INFO", user.ID, user.AvatarUrl, user.ProfileUrl)
	result, err := qtx.UpdateUser(ctx, postgresql.UpdateUserParams{
		ID: user.ID,
		AvatarUrl: sql.NullString{
			String: user.AvatarUrl,
			Valid:  user.AvatarUrl != "",
		},
		ProfileUrl: sql.NullString{
			String: user.ProfileUrl,
			Valid:  user.ProfileUrl != "",
		},
		Gender: sql.NullBool{
			Bool:  user.Gender,
			Valid: !user.Gender,
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Address: sql.NullString{
			String: user.Address,
			Valid:  user.Address != "",
		},
		DateOfBirth: sql.NullTime{
			Time:  user.DateOfBirth,
			Valid: !user.DateOfBirth.IsZero(),
		},
		Position: sql.NullString{
			String: user.Position,
			Valid:  user.Position != "",
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.UpdateUser failed")
	}
	return &domain.User{
		ID:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		FullName:    result.FullName.String,
		Role:        result.Role.String,
		AvatarUrl:   result.AvatarUrl.String,
		ProfileUrl:  result.ProfileUrl.String,
		Email:       result.Email,
		Password:    result.Password,
		Gender:      result.Gender.Bool,
		Phone:       result.Phone.String,
		Address:     result.Address.String,
		DateOfBirth: result.DateOfBirth.Time,
		Position:    result.Position.String,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, tx.Commit()
}

// GetUserIdsOfCompany implements auth.UserRepo.
func (rp *userRepo) GetUserIdsOfCompany(ctx context.Context, company string) ([]uuid.UUID, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	userIds, err := querier.GetUserIdsOfCompany(ctx, sql.NullString{
		String: company,
		Valid:  company != "",
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUserIdsOfCompany failed")
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
		NickName: sql.NullString{
			String: user.NickName,
			Valid:  user.NickName != "",
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.Create failed")
	}
	return &domain.User{
		ID:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		FullName:    result.FullName.String,
		Role:        result.Role.String,
		AvatarUrl:   result.AvatarUrl.String,
		ProfileUrl:  result.ProfileUrl.String,
		Email:       result.Email,
		Password:    result.Password,
		Gender:      result.Gender.Bool,
		Phone:       result.Phone.String,
		Address:     result.Address.String,
		DateOfBirth: result.DateOfBirth.Time,
		Position:    result.Position.String,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
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
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		FullName:    user.FullName.String,
		NickName:    user.FullName.String,
		Role:        user.Role.String,
		AvatarUrl:   user.AvatarUrl.String,
		ProfileUrl:  user.ProfileUrl.String,
		Password:    user.Password,
		Gender:      user.Gender.Bool,
		Phone:       user.Phone.String,
		Address:     user.Address.String,
		DateOfBirth: user.DateOfBirth.Time,
		Position:    user.Position.String,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// GetUsers implements auth.AuthRepo.
func (rp *userRepo) GetUsers(ctx context.Context, search string, limit, offset int32) ([]*domain.User, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	users, err := querier.GetUsers(ctx, postgresql.GetUsersParams{
		Column1: sql.NullString{
			String: search,
			Valid:  search != "",
		},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.GetUsers failed")
	}
	return lo.Map(users, func(user postgresql.AuthUser, _ int) *domain.User {
		return &domain.User{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName.String,
			NickName:    user.NickName.String,
			Role:        user.Role.String,
			AvatarUrl:   user.AvatarUrl.String,
			ProfileUrl:  user.ProfileUrl.String,
			Gender:      user.Gender.Bool,
			Phone:       user.Phone.String,
			Address:     user.Address.String,
			DateOfBirth: user.DateOfBirth.Time,
			Position:    user.Position.String,
			// Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}), nil
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
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		FullName:   user.FullName.String,
		NickName:   user.NickName.String,
		Role:       user.Role.String,
		AvatarUrl:  user.AvatarUrl.String,
		ProfileUrl: user.ProfileUrl.String,
		Password:   user.Password,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

var _ auth.UserRepo = (*userRepo)(nil)

func NewUserRepo(pg postgres.DBEngine) auth.UserRepo {
	return &userRepo{pg: pg}
}

var UserRepoSet = wire.NewSet(NewUserRepo)
