package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/search/config"
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
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
func (g *groupGRPCClient) GetGroups(ctx context.Context) (*v1.GetGroupsResponse, error) {
	client := v1.NewGroupServiceClient(g.conn)
	res, err := client.GetGroups(ctx, &v1.GetGroupsRequest{})
	if err != nil {
		slog.Warn("groupGRPCClient.GetGroupIdsByUserId failed", err)
		return nil, err
	}

	return res, nil

}
