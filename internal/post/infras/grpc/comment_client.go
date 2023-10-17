package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domain2 "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type commentGRPCClient struct {
	conn *grpc.ClientConn
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

// GetCommentsByPostID implements domain.CommentDomainService.
func (c *commentGRPCClient) GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error) {
	client := gen.NewCommentServiceClient(c.conn)

	res, err := client.GetCommentsByPostID(ctx, &gen.GetCommentsByPostIDRequest{
		PostID: postId.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentGRPCClient.GetCommentsByPostID failed")

	}
	results := make([]*sharedkernel.CommentHasChildren, 0)
	for _, item := range res.Comments {
		results = append(results, &sharedkernel.CommentHasChildren{
			ID:              uuid.Parse(item.Id),
			UserID:          item.UserId,
			Content:         item.Content,
			PostID:          item.PostId,
			ReplyToID:       item.ReplyToId,
			ParentCommentID: item.ParentCommentId,
			Children: lo.Map(item.Children, func(value *gen.CommentResponse, _ int) *domain2.Comment {
				return &domain2.Comment{
					ID:              value.Id,
					UserID:          uuid.Parse(value.UserId),
					Content:         value.Content,
					PostID:          uuid.Parse(value.PostId),
					ReplyToID:       uuid.Parse(value.ReplyToId),
					ParentCommentID: uuid.Parse(value.ParentCommentId),
					CreatedAt:       value.CreatedAt.AsTime(),
					UpdatedAt:       value.UpdatedAt.AsTime(),
				}
			}),
			CreatedAt: item.CreatedAt.AsTime(),
			UpdatedAt: item.UpdatedAt.AsTime(),
		})
	}
	return results, nil
}
