package repo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	UserID        = "20ead10d-2a4b-45d6-9bde-4f93fc76281d"
	UserEmail     = "wydqps@email.com"
	UserFirstName = "F"
	UserLastName  = "L"
	Company       = "@tlcmodular.com"
)

func TestCreateUser(t *testing.T) {
	repo := newUserRepo()
	arg := &domain.User{
		ID:        uuid.New(),
		Email:     utils.RandomEmail(),
		Password:  "123456789",
		FirstName: "F",
		LastName:  "L",
		FullName:  "FL",
	}
	user, err := repo.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.ID, user.ID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.FullName, user.FullName)
	require.NotZero(t, user.CreatedAt)
}
func TestGetUser(t *testing.T) {
	userFullName := sql.NullString{
		String: "FL",
		Valid:  true,
	}
	userId, err := uuid.Parse(UserID)
	require.NoError(t, err)
	repo := newUserRepo()
	user, err := repo.GetUser(context.Background(), userId)
	require.NoError(t, err)
	require.Equal(t, userId, user.ID)
	require.Equal(t, UserEmail, user.Email)
	require.Equal(t, UserFirstName, user.FirstName)
	require.Equal(t, UserLastName, user.LastName)
	require.Equal(t, userFullName, user.FullName)
	require.NotZero(t, user.CreatedAt)
}
func TestGetUserByEmail(t *testing.T) {
	userFullName := sql.NullString{
		String: "FL",
		Valid:  true,
	}
	userId, err := uuid.Parse(UserID)
	require.NoError(t, err)
	repo := newUserRepo()
	user, err := repo.GetUserByEmail(context.Background(), UserEmail)
	require.NoError(t, err)
	require.Equal(t, userId, user.ID)
	require.Equal(t, UserEmail, user.Email)
	require.Equal(t, UserFirstName, user.FirstName)
	require.Equal(t, UserLastName, user.LastName)
	require.Equal(t, userFullName, user.FullName)
	require.NotZero(t, user.CreatedAt)
}
func TestNewPostgresDB(t *testing.T) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)
	require.NotEmpty(t, cfg)
	db, err := postgres.NewPostgresDB(postgres.DBConnString(cfg.DbURL), postgres.DBConnReadString(cfg.DbRepURL))
	require.NoError(t, err)
	require.NotEmpty(t, db)
}

func newUserRepo() auth.UserRepo {
	cfg, _ := config.NewConfig()
	db, _ := postgres.NewPostgresDB(postgres.DBConnString(cfg.DbURL), postgres.DBConnReadString(cfg.DbRepURL))
	repo := NewUserRepo(db)
	return repo
}
