package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domain2 "github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type likeGRPCClient struct {
	conn *grpc.ClientConn
}

// GetLikesByPostID implements domain.LikeDomainService.
func (l *likeGRPCClient) GetLikesByPostID(ctx context.Context, postId uuid.UUID) ([]*domain2.Like, error) {
	panic("unimplemented")
}

var LikeGRPCClientSet = wire.NewSet(NewGRPCLikeClient)

var _ domain.LikeDomainService = (*likeGRPCClient)(nil)

func NewGRPCLikeClient(cfg *config.Config) (domain.LikeDomainService, error) {
	conn, err := grpc.Dial(cfg.LikeClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &likeGRPCClient{
		conn: conn,
	}, nil
}
