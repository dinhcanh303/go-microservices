package event

import (
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type UserCreated struct {
	sharedkernel.DomainEvent
	UserID       uuid.UUID `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	UserFullName string    `json:"user_full_name"`
	UserAvatar   string    `json:"user_avatar"`
}

func (e UserCreated) Identity() string {
	return "UserCreated"
}

type UserDeleted struct {
	sharedkernel.DomainEvent
	UserID uuid.UUID `json:"user_id" `
}

func (e UserDeleted) Identity() string {
	return "UserDeleted"
}

type GroupCreated struct {
	sharedkernel.DomainEvent
	GroupID     uuid.UUID `json:"group_id"`
	GroupName   string    `json:"group_name"`
	GroupAvatar string    `json:"group_avatar"`
}

func (e GroupCreated) Identity() string {
	return "GroupCreated"
}

type GroupDeleted struct {
	sharedkernel.DomainEvent
	GroupID uuid.UUID `json:"group_id"`
}

func (e GroupDeleted) Identity() string {
	return "GroupDeleted"
}
