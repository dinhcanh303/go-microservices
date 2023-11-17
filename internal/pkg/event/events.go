package event

import (
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type Notification struct {
	sharedkernel.DomainEvent
	ItemID      uuid.UUID `json:"item_id"`
	ItemType    string    `json:"item_type"`
	ItemContent string    `json:"item_content"`
}
