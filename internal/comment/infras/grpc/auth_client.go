package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGRPCClient struct {
	conn *grpc.ClientConn
}

// GetProfile implements domain.AuthDomainService.
func (a *authGRPCClient) GetProfile(ctx context.Context, id uuid.UUID) (*gen.GetProfileResponse, error) {
	client := gen.NewAuthServiceClient(a.conn)

	res, err := client.GetProfile(ctx, &gen.GetProfileRequest{
		Id: id.String(),
	})
	if err != nil {
		slog.Warn("authGRPCClient.GetProfile failed", err)
		return &gen.GetProfileResponse{}, err
	}
	return res, nil
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
