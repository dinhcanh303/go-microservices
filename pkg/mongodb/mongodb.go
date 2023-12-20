package mongodb

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnString string

type mongodb struct {
	db *mongo.Database
}

// Close implements MongoDBEngine.
func (mg *mongodb) Disconnect() {
	mg.db.Client().Disconnect(context.Background())
}

// GetCollection implements MongoDBEngine.
func (mg *mongodb) GetCollection(collectionName string) *mongo.Collection {
	return mg.db.Collection(collectionName)
}

// Connect implements MongoDB.
func NewMongoDB(connectionString MongoDBConnString, dbName string) (MongoDBEngine, error) {
	mg := &mongodb{}
	clientOptions := options.Client().ApplyURI(string(connectionString))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	mg.db = client.Database(dbName)
	slog.Info("📰 connected to mongodb 🎉")
	return mg, nil
}

var _ MongoDBEngine = (*mongodb)(nil)
