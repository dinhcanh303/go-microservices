package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.AuthDomainService = (*authGRPCClient)(nil)

var AuthGRPCClientSet = wire.NewSet(NewGRPCAuthClient)

func NewGRPCAuthClient(cfg *config.Config) (domain.AuthDomainService, error) {
	conn, err := grpc.Dial(cfg.AuthClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &authGRPCClient{
		conn: conn,
	}, nil
}

// GetAllUserIdByUserId implements domain.AuthDomainService.
func (a *authGRPCClient) GetUsers(ctx context.Context) (*gen.GetUsersResponse, error) {
	client := gen.NewAuthServiceClient(a.conn)
	res, err := client.GetUsers(ctx, &gen.GetUsersRequest{})
	if err != nil {
		slog.Warn("authGRPCClient.GetUsers failed", err)
		return nil, err
	}

	return res, nil
}
