package keys

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/google/uuid"
)

type KeyRepo interface {
	CreateKey(context.Context, *domain.Key) (*domain.Key, error)
	UpdateKeyByUserID(context.Context, *domain.Key) (*domain.Key, error)
	FindKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.Key, error)
	DeleteKeyByID(ctx context.Context, id int64) error
	DeleteKeyByUserID(ctx context.Context, userID uuid.UUID) error
	FindKeyByRefreshToken(ctx context.Context, refreshToken string) (*domain.Key, error)
	FindKeyByRefreshTokenUsed(ctx context.Context, refreshToken string) (*domain.Key, error)
}
type UseCase interface {
	CreateKeyToken(context.Context, *domain.Key) (*domain.Key, error)
	FindKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.Key, error)
	DeleteKeyByID(ctx context.Context, id int64) error
	DeleteKeyByUserID(ctx context.Context, userID uuid.UUID) error
	FindKeyByRefreshToken(ctx context.Context, refreshToken string) (*domain.Key, error)
	FindKeyByRefreshTokenUsed(ctx context.Context, refreshToken string) (*domain.Key, error)
}
