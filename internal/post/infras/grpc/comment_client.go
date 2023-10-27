package grpc

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainComment "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
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
			ID:              uuid.MustParse(item.Id),
			UserID:          uuid.MustParse(item.UserId),
			Content:         item.Content,
			PostID:          uuid.MustParse(item.PostId),
			ReplyToID:       utils.StringToNullUUID(item.ReplyToId),
			ParentCommentID: utils.StringToNullUUID(item.ParentCommentId),
			Children: lo.Map(item.Children, func(value *gen.CommentResponseHasLike, _ int) *domainComment.CommentHasLike {
				return &domainComment.CommentHasLike{
					ID:      uuid.MustParse(value.Id),
					UserID:  uuid.MustParse(value.UserId),
					Content: value.Content,
					Likes: lo.Map(item.Likes, func(value *gen.LikeResponseInComment, _ int) *domainLike.Like {
						return &domainLike.Like{
							ID:           uuid.MustParse(value.Id),
							UserID:       uuid.MustParse(value.UserId),
							Emoji:        value.Emoji,
							LikeableType: value.LikeableType,
							LikeableID:   uuid.MustParse(value.LikeableId),
							CreatedAt:    value.CreatedAt.AsTime(),
							UpdatedAt:    value.UpdatedAt.AsTime(),
						}
					}),
					PostID:          uuid.MustParse(value.PostId),
					ReplyToID:       utils.StringToNullUUID(value.ReplyToId),
					ParentCommentID: utils.StringToNullUUID(value.ParentCommentId),
					CreatedAt:       value.CreatedAt.AsTime(),
					UpdatedAt:       value.UpdatedAt.AsTime(),
				}
			}),
			Likes: lo.Map(item.Likes, func(value *gen.LikeResponseInComment, _ int) *domainLike.Like {
				return &domainLike.Like{
					ID:           uuid.MustParse(value.Id),
					UserID:       uuid.MustParse(value.UserId),
					Emoji:        value.Emoji,
					LikeableType: value.LikeableType,
					LikeableID:   uuid.MustParse(value.LikeableId),
					CreatedAt:    value.CreatedAt.AsTime(),
					UpdatedAt:    value.UpdatedAt.AsTime(),
				}
			}),
			CreatedAt: item.CreatedAt.AsTime(),
			UpdatedAt: item.UpdatedAt.AsTime(),
		})
	}
	return results, nil
}
