package token_test

import (
	"testing"
	"time"

	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	maker := token.NewJWTMaker()
	id := utils.RandomUUID()
	email := utils.RandomEmailCompany()
	fullName := utils.RandomString(8)
	avatarUrl := utils.RandomString(10)
	role := utils.User
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	publicKey, err := utils.GenerateRandomHexBytes(64)
	require.NoError(t, err)
	require.NotEmpty(t, publicKey)
	token, err := maker.CreateToken(&token.Payload{
		ID:        id,
		Email:     email,
		FullName:  fullName,
		Role:      role,
		AvatarUrl: avatarUrl,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, publicKey)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token, publicKey)
	require.NoError(t, err)
	require.Equal(t, id, payload.ID)
	require.Equal(t, email, payload.Email)
	require.Equal(t, fullName, payload.FullName)
	require.Equal(t, avatarUrl, payload.AvatarUrl)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
