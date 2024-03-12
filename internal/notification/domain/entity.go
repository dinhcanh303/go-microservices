package domain

import (
	"time"
)

type Notification struct {
	ID         int64                  `json:"id" gorm:"column:id"`
	ActorID    string                 `json:"actor_id" gorm:"column:actor_id"`
	SenderID   string                 `json:"sender_id" gorm:"column:sender_id"`
	Data       map[string]interface{} `json:"data" gorm:"serializer:json"`
	Type       string                 `json:"type" gorm:"column:type"`
	ObjectType string                 `json:"object_type" gorm:"column:object_type"`
	ObjectID   string                 `json:"object_id" gorm:"column:object_id"`
	ReadAt     *time.Time             `json:"read_at" gorm:"column:read_at"`
	CreatedAt  time.Time              `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time              `json:"updated_at" gorm:"column:updated_at"`
}

func (Notification) TableName() string {
	return "noti.notifications"
}
