//go:build wireinject
// +build wireinject

package app

import (
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/router"
	"github.com/dinhcanh303/go-microservices/internal/auth/infras/repo"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	cfgLdap *configs.Ldap,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		ldapClientFunc,
		router.AuthGRPCServerSet,
		auth.UseCaseSet,
		keys.UseCaseSet,
		repo.KeyRepoSet,
		repo.UserRepoSet,
	))
}
func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
func ldapClientFunc(config *configs.Ldap) (ldap.LdapClient, func(), error) {
	ldapClient, err := ldap.NewLdapClient(config, []string{""})
	if err != nil {
		slog.Warn("Connect failed:", err)
		return nil, nil, err
	}
	return ldapClient, func() { ldapClient.Close() }, nil
}
