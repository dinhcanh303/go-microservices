package repo

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	"github.com/dinhcanh303/go-microservices/pkg/mongodb"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type notificationRepo struct {
	mongo mongodb.MongoDBEngine
}

const collectionName = "notifications"

var RepositorySet = wire.NewSet()

func NewNotificationRepo(mongo mongodb.MongoDBEngine) notifications.NotificationRepo {
	return &notificationRepo{
		mongo: mongo,
	}
}

// GetAllByID implements notifications.NotificationRepo.
func (n *notificationRepo) GetAllByID(context.Context, uuid.UUID) ([]*domain.Notification, error) {
	panic("unimplemented")
}

// GetListByID implements notifications.NotificationRepo.
func (n *notificationRepo) GetListByID(ctx context.Context, id uuid.UUID, offset uint32, limit uint32) ([]*domain.Notification, error) {
	// collection := n.mongo.GetCollection(collectionName)
	panic("unimplemented")
}

// Update implements notifications.NotificationRepo.
func (n *notificationRepo) Update(ctx context.Context, noti *domain.Notification) error {
	collection := n.mongo.GetCollection(collectionName)
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "key", Value: noti.Key}}
	update := bson.D{{"$set", bson.D{{Key: ""}}}}
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

var _ notifications.NotificationRepo = (*notificationRepo)(nil)
