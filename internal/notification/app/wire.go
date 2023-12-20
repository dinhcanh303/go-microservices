///go:build wireinject
/// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/notification/config"
	"github.com/dinhcanh303/go-microservices/pkg/mongodb"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr mongodb.MongoDBEngine,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
	))
}
func dbEngineFunc(url mongodb.MongoDBConnString, dbName string) (mongodb.MongoDBEngine, func(), error) {
	db, err := mongodb.NewMongoDB(url, dbName)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Disconnect() }, nil
}
