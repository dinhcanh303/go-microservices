package domain

import (
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type GroupSearch struct {
	sharedkernel.AggregateRoot
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Type   string    `json:"type"`
}

func NewGroupSearch(e event.GroupCreated) *GroupSearch {
	groupSearch := &GroupSearch{
		ID:     e.ID,
		Name:   e.Name,
		Avatar: e.Avatar,
		Type:   "group",
	}
	return groupSearch
}
