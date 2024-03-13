package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type (
	AuthDomainService interface {
		GetProfile(ctx context.Context, id string) (*gen.GetProfileResponse, error)
	}
)
