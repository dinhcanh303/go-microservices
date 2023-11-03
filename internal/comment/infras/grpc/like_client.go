package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type likeGRPCClient struct {
	conn *grpc.ClientConn
}

// GetLikesByPostID implements domain.LikeDomainService.
func (l *likeGRPCClient) GetLikesByCommentID(ctx context.Context, commentId uuid.UUID) ([]*domainLike.Like, error) {
	client := gen.NewLikeServiceClient(l.conn)
	res, err := client.GetLikesByCommentID(ctx, &gen.GetLikesByCommentIDRequest{
		CommentId: commentId.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentGRPCClient.GetCommentsByPostID failed")

	}
	results := make([]*domainLike.Like, 0)
	for _, item := range res.Likes {
		results = append(results, &domainLike.Like{
			ID:           uuid.MustParse(item.Id),
			UserID:       uuid.MustParse(item.UserId),
			Emoji:        item.Emoji,
			LikeableType: item.LikeableType,
			LikeableID:   uuid.MustParse(item.LikeableId),
			CreatedAt:    item.CreatedAt.AsTime(),
			UpdatedAt:    item.UpdatedAt.AsTime(),
		})
	}
	return results, nil
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
