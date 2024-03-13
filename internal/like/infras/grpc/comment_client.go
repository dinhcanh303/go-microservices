package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type commentGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.CommentDomainService = (*commentGRPCClient)(nil)

var CommentGRPCClientSet = wire.NewSet(NewGRPCCommentClient)

func NewGRPCCommentClient(cfg *config.Config) (domain.CommentDomainService, error) {
	conn, err := grpc.Dial(cfg.CommentClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &commentGRPCClient{
		conn: conn,
	}, nil
}

// GetComment implements domain.CommentDomainService.
func (c *commentGRPCClient) GetComment(ctx context.Context, id uuid.UUID) (*gen.GetCommentResponse, error) {
	client := gen.NewCommentServiceClient(c.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetComment(ctxBackground, &gen.GetCommentRequest{
		Id: id.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "postGRPCClient.GetPostNormal failed")
	}
	return res, nil
}
