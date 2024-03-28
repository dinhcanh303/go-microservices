package domain

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	"github.com/google/uuid"
)

type AuthDomainService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*v1.GetProfileResponse, error)
}
