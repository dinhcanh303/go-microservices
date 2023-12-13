package domain

import (
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type UserSearch struct {
	sharedkernel.AggregateRoot
	UserID       uuid.UUID
	UserFullName string
	UserEmail    string
	UserAvatar   string
}

func NewUserSearch(e event.UserCreated) *UserSearch {
	userSearch := &UserSearch{
		UserID:       e.UserID,
		UserFullName: e.UserFullName,
		UserEmail:    e.UserEmail,
		UserAvatar:   e.UserAvatar,
	}
	return userSearch
}
