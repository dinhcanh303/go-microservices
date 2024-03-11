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

type PostNoti struct {
	sharedkernel.DomainEvent
	ActorID    string                 `json:"actor_id"`
	SenderIDs  []string               `json:"sender_ids"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	ObjectType string                 `json:"object_type"`
	ObjectID   string                 `json:"object_id"`
}

func (e *PostNoti) Identity() string {
	return "PostNoti"
}

type CommentNoti struct {
	sharedkernel.DomainEvent
	ActorID    string                 `json:"actor_id"`
	SenderIDs  []string               `json:"sender_ids"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	ObjectType string                 `json:"object_type"`
	ObjectID   string                 `json:"object_id"`
}

func (e *CommentNoti) Identity() string {
	return "CommentNoti"
}

type LikeNoti struct {
	sharedkernel.DomainEvent
	ActorID    string                 `json:"actor_id"`
	SenderIDs  []string               `json:"sender_ids"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	ObjectType string                 `json:"object_type"`
	ObjectID   string                 `json:"object_id"`
}

func (e *LikeNoti) Identity() string {
	return "LikeNoti"
}
