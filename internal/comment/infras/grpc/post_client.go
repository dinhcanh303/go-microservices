package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type postGRPCClient struct {
	conn *grpc.ClientConn
}

var PostGRPCClientSet = wire.NewSet(NewGRPCPostClient)

var _ domain.PostDomainService = (*postGRPCClient)(nil)

func NewGRPCPostClient(cfg *config.Config) (domain.PostDomainService, error) {
	conn, err := grpc.Dial(cfg.PostClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &postGRPCClient{
		conn: conn,
	}, nil
}

// GetLikesByPostID implements domain.LikeDomainService.
func (l *postGRPCClient) GetPostNormal(ctx context.Context, id uuid.UUID) (*gen.GetPostNormalResponse, error) {
	client := gen.NewPostServiceClient(l.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetPostNormal(ctxBackground, &gen.GetPostNormalRequest{
		Id: id.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "postGRPCClient.GetPostNormal failed")
	}
	return res, nil
}
