package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainComment "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type commentGRPCClient struct {
	conn *grpc.ClientConn
}

// CountCommentByPostID implements domain.CommentDomainService.
func (c *commentGRPCClient) CountCommentByPostID(ctx context.Context, postId uuid.UUID) (int64, error) {
	client := gen.NewCommentServiceClient(c.conn)
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return 0, err
	}
	res, err := client.CountCommentByPostID(ctxBackground, &gen.CountCommentByPostIDRequest{
		PostId: postId.String(),
	})
	if err != nil {
		slog.Warn("commentGRPCClient.GetCommentsByPostID failed", err)
		return 0, nil
	}
	return res.Count, nil
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
	ctxBackground, err := utils.OutgoingContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := client.GetCommentsByPostID(ctxBackground, &gen.GetCommentsByPostIDRequest{
		PostId: postId.String(),
	})
	results := make([]*sharedkernel.CommentHasChildren, 0)
	if err != nil {
		slog.Warn("commentGRPCClient.GetCommentsByPostID failed", err)
		return results, nil
	}
	for _, item := range res.Comments {
		results = append(results, &sharedkernel.CommentHasChildren{
			ID:              uuid.MustParse(item.Id),
			UserID:          uuid.MustParse(item.UserId),
			Content:         item.Content,
			PostID:          uuid.MustParse(item.PostId),
			ReplyToID:       utils.StringToNullUUIDNormal(item.ReplyToId),
			ParentCommentID: utils.StringToNullUUIDNormal(item.ParentCommentId),
			Children: lo.Map(item.Children, func(value *gen.CommentHasMetadata, _ int) *domainComment.CommentHasMetadata {
				return &domainComment.CommentHasMetadata{
					ID:      uuid.MustParse(value.Id),
					UserID:  uuid.MustParse(value.UserId),
					Content: value.Content,
					Likes: &domainLike.LikesInfo{
						YourLikedEmoji:    item.Likes.YourLikedEmoji,
						YourLike:          item.Likes.YourLike,
						OthersLikedEmojis: item.Likes.OthersLikedEmojis,
						OthersLikes:       item.Likes.OthersLikes,
					},
					PostID:          uuid.MustParse(value.PostId),
					ReplyToID:       utils.StringToNullUUIDNormal(value.ReplyToId),
					ParentCommentID: utils.StringToNullUUIDNormal(value.ParentCommentId),
					CreatedAt:       value.CreatedAt.AsTime(),
					UpdatedAt:       value.UpdatedAt.AsTime(),
				}
			}),
			Likes: &domainLike.LikesInfo{
				YourLikedEmoji:    item.Likes.YourLikedEmoji,
				YourLike:          item.Likes.YourLike,
				OthersLikedEmojis: item.Likes.OthersLikedEmojis,
				OthersLikes:       item.Likes.OthersLikes,
			},
			CreatedAt: item.CreatedAt.AsTime(),
			UpdatedAt: item.UpdatedAt.AsTime(),
		})
	}
	return results, nil
}
