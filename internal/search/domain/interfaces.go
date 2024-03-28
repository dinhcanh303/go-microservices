package domain

import (
	"context"

	v1a "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1g "github.com/dinhcanh303/go-microservices/api/group/v1"
)

type (
	GroupDomainService interface {
		GetGroups(ctx context.Context) (*v1g.GetGroupsResponse, error)
	}
	AuthDomainService interface {
		GetUsers(ctx context.Context) (*v1a.GetUsersResponse, error)
	}
)
