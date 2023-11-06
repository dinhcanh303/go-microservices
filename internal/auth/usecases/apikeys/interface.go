package apikeys

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
)

type ApiKeyRepo interface {
	CreateApiKey(context.Context, *domain.ApiKey) (*domain.ApiKey, error)
}
type UseCase interface {
	CreateApiKey(context.Context, *domain.ApiKey) (*domain.ApiKey, error)
}
