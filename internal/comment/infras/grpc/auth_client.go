package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/google/uuid"
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

// GetProfile implements domain.AuthDomainService.
func (a *authGRPCClient) GetProfile(ctx context.Context, id uuid.UUID) (*v1.GetProfileResponse, error) {
	client := v1.NewAuthServiceClient(a.conn)

	res, err := client.GetProfile(ctx, &v1.GetProfileRequest{
		Id: id.String(),
	})
	if err != nil {
		slog.Warn("authGRPCClient.GetProfile failed", err)
		return &v1.GetProfileResponse{}, err
	}
	return res, nil
}
