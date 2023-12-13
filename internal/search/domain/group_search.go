package domain

import (
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type GroupSearch struct {
	sharedkernel.AggregateRoot
	GroupID     uuid.UUID
	GroupName   string
	GroupAvatar string
}

func NewGroupSearch(e event.GroupCreated) *GroupSearch {
	groupSearch := &GroupSearch{
		GroupID:     e.GroupID,
		GroupName:   e.GroupName,
		GroupAvatar: e.GroupAvatar,
	}
	return groupSearch
}
