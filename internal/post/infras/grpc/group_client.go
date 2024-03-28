package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
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

// GetGroupMembers implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroupMembers(ctx context.Context, groupId uuid.NullUUID) (*v1.GetGroupMembersResponse, error) {
	client := v1.NewGroupServiceClient(g.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetGroupMembers(ctxBackground, &v1.GetGroupMembersRequest{
		GroupId: groupId.UUID.String(),
	})
	if err != nil {
		slog.Warn("groupGRPCClient.GetGroup failed", err)
		return nil, err
	}
	return res, nil
}

// GetGroup implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroup(ctx context.Context, groupId uuid.NullUUID) (*v1.GetGroupResponse, error) {
	client := v1.NewGroupServiceClient(g.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetGroup(ctxBackground, &v1.GetGroupRequest{
		Id: groupId.UUID.String(),
	})
	if err != nil {
		slog.Warn("groupGRPCClient.GetGroup failed", err)
		return nil, err
	}
	return res, nil
}

// GetAllGroupIdByUserId implements domain.GroupDomainService.
func (g *groupGRPCClient) GetGroupIdsByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	client := v1.NewGroupServiceClient(g.conn)
	res, err := client.GetGroupIdsByUserId(ctx, &v1.GetGroupIdsByUserIdRequest{
		UserId: userId.String(),
	})
	results := make([]uuid.UUID, 0)
	if err != nil {
		slog.Warn("groupGRPCClient.GetGroupIdsByUserId failed", err)
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
