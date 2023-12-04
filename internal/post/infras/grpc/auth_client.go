package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGRPCClient struct {
	conn *grpc.ClientConn
}

// GetAllUserIdByUserId implements domain.AuthDomainService.
func (a *authGRPCClient) GetAllUserIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	client := gen.NewAuthServiceClient(a.conn)

	res, err := client.GetAllUserIdByUserId(ctx, &gen.GetAllUserIdByUserIdRequest{
		UserId: userId.String(),
	})
	results := make([]uuid.UUID, 0)
	if err != nil {
		slog.Warn("authGRPCClient.GetAllUserIdByUserId failed", err)
		return results, nil
	}
	for _, item := range res.UserIds {
		uuid, err := uuid.Parse(item)
		if err == nil {
			results = append(results, uuid)
		}
	}
	return results, nil
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
