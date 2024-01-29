package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
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

// GetLikesByPostID implements domain.LikeDomainService.
func (l *likeGRPCClient) GetLikesInfoByCommentID(ctx context.Context, commentId, userId uuid.UUID) (*domainLike.LikesInfo, error) {
	client := gen.NewLikeServiceClient(l.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetLikesInfoByCommentID(ctxBackground, &gen.GetLikesInfoByCommentIDRequest{
		CommentId: commentId.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "commentGRPCClient.GetLikesInfoByCommentID failed")

	}

	return &domainLike.LikesInfo{
		YourLikedEmoji:    res.Likes.YourLikedEmoji,
		YourLike:          res.Likes.YourLike,
		OthersLikedEmojis: res.Likes.OthersLikedEmojis,
		OthersLikes:       res.Likes.OthersLikes,
	}, nil
}
