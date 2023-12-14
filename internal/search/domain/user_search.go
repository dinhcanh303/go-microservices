package domain

import (
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type UserSearch struct {
	sharedkernel.AggregateRoot
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
	Type   string    `json:"type"`
}

func NewUserSearch(e event.UserCreated) *UserSearch {
	userSearch := &UserSearch{
		ID:     e.ID,
		Name:   e.Name,
		Email:  e.Email,
		Avatar: e.Avatar,
		Type:   "user",
	}
	return userSearch
}
