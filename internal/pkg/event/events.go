package event

import (
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type UserCreated struct {
	sharedkernel.DomainEvent
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Type   string    `json:"type"`
}

func (e UserCreated) Identity() string {
	return "UserCreated"
}

type UserDeleted struct {
	sharedkernel.DomainEvent
	ID uuid.UUID `json:"id" `
}

func (e UserDeleted) Identity() string {
	return "UserDeleted"
}

type GroupCreated struct {
	sharedkernel.DomainEvent
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Type   string    `json:"type"`
}

func (e GroupCreated) Identity() string {
	return "GroupCreated"
}

type GroupDeleted struct {
	sharedkernel.DomainEvent
	ID uuid.UUID `json:"id"`
}

func (e GroupDeleted) Identity() string {
	return "GroupDeleted"
}

type Notification struct {
	sharedkernel.DomainEvent
	Key      string               `bson:"key,omitempty,unique" json:"key"`
	Subject  sharedkernel.Subject `bson:"subject" json:"subject"`
	DiObject sharedkernel.Subject `bson:"di_object" json:"di_object"`
	InObject sharedkernel.Subject `bson:"in_object" json:"in_object"`
	PrObject sharedkernel.Subject `bson:"pr_object" json:"pr_object"`
}

func (e *Notification) Identity() string {
	return "Notification"
}
