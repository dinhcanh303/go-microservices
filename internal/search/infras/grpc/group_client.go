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

type groupGRPCClient struct {
	conn *grpc.ClientConn
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

// GetAllGroupIdByUserId implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroups(ctx context.Context) (*gen.GetGroupsResponse, error) {
	client := gen.NewGroupServiceClient(g.conn)
	res, err := client.GetGroups(ctx, &gen.GetGroupsRequest{})
	if err != nil {
		slog.Warn("groupGRPCClient.GetGroupIdsByUserId failed", err)
		return nil, err
	}

	return res, nil

}
