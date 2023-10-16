package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type commentGRPCClient struct {
	conn *grpc.ClientConn
}

// GetCommentsByPostID implements domain.CommentDomainService.
func (*commentGRPCClient) GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.CommentItem, error) {
	panic("unimplemented")
}

var CommentGRPCClientSet = wire.NewSet(NewGRPCCommentClient)

var _ domain.CommentDomainService = (*commentGRPCClient)(nil)

func NewGRPCCommentClient(cfg *config.Config) (domain.CommentDomainService, error) {
	conn, err := grpc.Dial(cfg.CommentClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &commentGRPCClient{
		conn: conn,
	}, nil
}
