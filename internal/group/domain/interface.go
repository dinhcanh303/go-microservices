package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
)

type AuthDomainService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*gen.GetProfileResponse, error)
}
