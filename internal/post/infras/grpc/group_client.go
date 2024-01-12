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

type groupGRPCClient struct {
	conn *grpc.ClientConn
}

// GetAllGroupIdByUserId implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroupIdsByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	client := gen.NewGroupServiceClient(g.conn)

	res, err := client.GetGroupIdsByUserId(ctx, &gen.GetGroupIdsByUserIdRequest{
		UserId: userId.String(),
	})
	results := make([]uuid.UUID, 0)
	if err != nil {
		slog.Warn("commentGRPCClient.GetGroupIdsByUserId failed", err)
		return results, nil
	}
	for _, item := range res.GroupIds {
		uuid, err := uuid.Parse(item)
		if err == nil {
			results = append(results, uuid)
		}
	}
	return results, nil

}

var _ domain.GroupDomainService = (*groupGRPCClient)(nil)

var GroupGRPCClientSet = wire.NewSet(NewGRPCGroupClient)

func NewGRPCGroupClient(cfg *config.Config) (domain.GroupDomainService, error) {
	conn, err := grpc.Dial(cfg.GroupClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &groupGRPCClient{
		conn: conn,
	}, nil
}
