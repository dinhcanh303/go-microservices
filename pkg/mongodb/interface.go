package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type MongoDBEngine interface {
	GetCollection(collectionName string) *mongo.Collection
	Disconnect()
}
