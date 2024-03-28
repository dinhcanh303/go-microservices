package grpc

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type groupGRPCClient struct {
	conn *grpc.ClientConn
}

// GetGroupMembers implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroupMembers(ctx context.Context, groupId string) (*v1.GetGroupMembersResponse, error) {
	client := v1.NewGroupServiceClient(g.conn)
	res, err := client.GetGroupMembers(ctx, &v1.GetGroupMembersRequest{
		GroupId: groupId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil

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
