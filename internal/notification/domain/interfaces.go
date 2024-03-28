package domain

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
)

type (
	AuthDomainService interface {
		GetProfile(ctx context.Context, id string) (*v1.GetProfileResponse, error)
	}
)
