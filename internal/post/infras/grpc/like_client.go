package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/like/v1"
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type likeGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.LikeDomainService = (*likeGRPCClient)(nil)

var LikeGRPCClientSet = wire.NewSet(NewGRPCLikeClient)

func NewGRPCLikeClient(cfg *config.Config) (domain.LikeDomainService, error) {
	conn, err := grpc.Dial(cfg.LikeClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &likeGRPCClient{
		conn: conn,
	}, nil
}

// GetLikesByPostID implements domain.LikeDomainService.
func (l *likeGRPCClient) GetLikesByPostID(ctx context.Context, postId uuid.UUID) (*domainLike.LikesInfo, error) {
	client := v1.NewLikeServiceClient(l.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetLikesInfoByPostID(ctxBackground, &v1.GetLikesInfoByPostIDRequest{
		PostId: postId.String(),
	})
	results := &domainLike.LikesInfo{}
	if err != nil {
		slog.Warn("likeGRPCClient.GetLikesByPostID failed", err)
		return results, nil
	}
	return &domainLike.LikesInfo{
		YourLikedEmoji:    res.Likes.YourLikedEmoji,
		YourLike:          res.Likes.YourLike,
		OthersLikedEmojis: res.Likes.OthersLikedEmojis,
		OthersLikes:       res.Likes.OthersLikes,
	}, nil
}
