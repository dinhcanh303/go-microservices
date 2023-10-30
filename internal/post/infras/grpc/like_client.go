package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type likeGRPCClient struct {
	conn *grpc.ClientConn
}

// GetLikesByPostID implements domain.LikeDomainService.
func (l *likeGRPCClient) GetLikesByPostID(ctx context.Context, postId uuid.UUID) ([]*domainLike.Like, error) {
	client := gen.NewLikeServiceClient(l.conn)
	res, err := client.GetLikesByPostID(ctx, &gen.GetLikesByPostIDRequest{
		PostId: postId.String(),
	})
	results := make([]*domainLike.Like, 0)
	if err != nil {
		slog.Warn("likeGRPCClient.GetLikesByPostID failed", err)
		return results, nil
	}
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
