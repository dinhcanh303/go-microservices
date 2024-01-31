package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type (
	GroupDomainService interface {
		GetGroups(ctx context.Context) (*gen.GetGroupsResponse, error)
	}
	AuthDomainService interface {
		GetUsers(ctx context.Context) (*gen.GetUsersResponse, error)
	}
)
