package repo

import (
	"context"
	"testing"

	"github.com/dinhcanh303/go-microservices/cmd/gateway/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestKeyPostgres(t *testing.T) {
	publicKey, err := utils.GenerateRandomHexBytes(64)
	require.NoError(t, err)
	require.NotEmpty(t, publicKey)
	privateKey, err := utils.GenerateRandomHexBytes(64)
	require.NoError(t, err)
	require.NotEmpty(t, privateKey)
	userID := uuid.New()
	k1 := domain.NewKey(userID, publicKey, privateKey)
	cfg, err := config.NewConfig()
	require.NoError(t, err)
	require.NotEmpty(t, cfg)
	db, err := postgres.NewPostgresDB(postgres.DBConnString(cfg.DsnURL))
	require.NoError(t, err)
	require.NotEmpty(t, db)
	repo := NewKeyRepo(db)
	// Test CreateKey
	key, err := repo.CreateKey(context.Background(), k1)
	require.NoError(t, err)
	require.Equal(t, key.UserID, userID)
	require.Equal(t, key.PublicKey, publicKey)
	require.Equal(t, key.PrivateKey, privateKey)

	// Test FindKeyByUserID
	key, err = repo.FindKeyByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, key.UserID, userID)
	require.Equal(t, key.PublicKey, publicKey)
	require.Equal(t, key.PrivateKey, privateKey)
}
