package apikeys

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
)

type service struct {
	repo ApiKeyRepo
}

// CreateApiKey implements UseCase.
func (s *service) CreateApiKey(context.Context, *domain.ApiKey) (*domain.ApiKey, error) {

	panic("unimplemented")
}

var _ UseCase = (*service)(nil)

func NewService(repo ApiKeyRepo) UseCase {
	return &service{
		repo: repo,
	}
}
