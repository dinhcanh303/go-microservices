package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type commentGRPCServer struct {
	gen.UnimplementedCommentServiceServer
	cfg *config.Config
	uc  comments.UseCase
}

var _ gen.CommentServiceServer = (*commentGRPCServer)(nil)

var CommentGRPCServerSet = wire.NewSet(NewGRPCCommentServer)

func NewGRPCCommentServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc comments.UseCase,
) gen.CommentServiceServer {
	svc := commentGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	gen.RegisterCommentServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

// CountCommentByCommentID implements gen.CommentServiceServer.
func (c *commentGRPCServer) CountCommentByCommentID(ctx context.Context, request *gen.CountCommentByCommentIDRequest) (*gen.CountCommentByCommentIDResponse, error) {
	slog.Info("POST: CountCommentByCommentID")
	commentId, err := uuid.Parse(request.CommentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	count, err := c.uc.CountCommentByCommentID(ctx, commentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count")
	}
	return &gen.CountCommentByCommentIDResponse{
		Count: count,
	}, nil
}

// CountCommentByPostID implements gen.CommentServiceServer.
func (c *commentGRPCServer) CountCommentByPostID(ctx context.Context, request *gen.CountCommentByPostIDRequest) (*gen.CountCommentByPostIDResponse, error) {
	slog.Info("POST: CountCommentByPostID")
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	count, err := c.uc.CountCommentByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count")
	}
	return &gen.CountCommentByPostIDResponse{
		Count: count,
	}, nil
}

// CreateComment implements gen.CommentServiceServer.
func (c *commentGRPCServer) CreateComment(ctx context.Context, request *gen.CreateCommentRequest) (*gen.CreateCommentResponse, error) {
	slog.Info("POST: CreateComment")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	postId, err := uuid.Parse(request.Comment.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	parentCommentId, _ := uuid.Parse(request.Comment.ParentCommentId)
	tagIds, _ := utils.ConvertArStringToArUUID(request.Comment.TagIds)
	replyId, _ := uuid.Parse(request.Comment.ReplyId)
	model := domain.Comment{
		ID:      uuid.New(),
		PostID:  postId,
		Content: request.Comment.Content,
		ParentCommentID: uuid.NullUUID{
			UUID:  parentCommentId,
			Valid: request.Comment.ParentCommentId != "",
		},
		ReplyID: uuid.NullUUID{
			UUID:  replyId,
			Valid: request.Comment.ReplyId != "",
		},
		TagIDs: tagIds,
		UserID: user.ID,
	}
	comment, err := c.uc.CreateComment(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreateComment failed")
	}
	return &gen.CreateCommentResponse{
		Comment: &gen.Comment{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
			Content:         comment.Content,
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			UserId:          comment.UserID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
		},
	}, nil
}

// DeleteComment implements gen.CommentServiceServer.
func (c *commentGRPCServer) DeleteComment(ctx context.Context, request *gen.DeleteCommentRequest) (*gen.DeleteCommentResponse, error) {
	slog.Info("DELETE: DeleteComment")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	result, err := c.uc.DeleteComment(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete")
	}
	return &gen.DeleteCommentResponse{
		Deleted: result,
	}, nil
}

// GetComment implements gen.CommentServiceServer.
func (c *commentGRPCServer) GetComment(ctx context.Context, request *gen.GetCommentRequest) (*gen.GetCommentResponse, error) {
	slog.Info("GET: GetComment")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	comment, err := c.uc.GetComment(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comment")
	}
	return &gen.GetCommentResponse{
		Comment: &gen.Comment{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			Content:         comment.Content,
			UserId:          comment.UserID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
		},
	}, nil
}
func (c *commentGRPCServer) GetCommentsByCommentID(ctx context.Context, request *gen.GetCommentsByCommentIDRequest) (*gen.GetCommentsByCommentIDResponse, error) {
	slog.Info("GET: GetCommentsByCommentID")
	commentId, err := uuid.Parse(request.CommentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	res := gen.GetCommentsByCommentIDResponse{}
	if request.Limit == 0 {
		request.Limit = 10
	}
	comments, err := c.uc.GetCommentsByCommentID(ctx, commentId, user.ID, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments by post ID")
	}
	for _, comment := range comments {
		res.Comments = append(res.Comments, &gen.CommentHasMetadata{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			UserId:          comment.UserID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
			Content:         comment.Content,
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
			Likes: &gen.LikeInfo{
				YourLikedEmoji:    comment.Likes.YourLikedEmoji,
				YourLike:          comment.Likes.YourLike,
				OthersLikedEmojis: comment.Likes.OthersLikedEmojis,
				OthersLikes:       comment.Likes.OthersLikes,
			},
			Attachments: lo.Map(comment.Attachments, func(item *domainUpload.Attachment, _ int) *gen.Attachment {
				return &gen.Attachment{
					Id:             item.ID.String(),
					UserId:         item.UserID.String(),
					AttachableType: item.AttachableType,
					AttachableId:   item.AttachableID.String(),
					Filename:       item.FileName,
					Url:            item.URL,
					UrlThumbnail:   item.URLThumbnail,
					Extension:      item.Extension,
					MimeType:       item.MimeType,
					Folder:         item.Folder,
					CreatedAt:      timestamppb.New(item.CreatedAt),
					UpdatedAt:      timestamppb.New(item.UpdatedAt),
				}
			}),
		})
	}
	return &res, nil
}

// ListCommentByPostID implements gen.CommentServiceServer.
func (c *commentGRPCServer) GetCommentsByPostID(ctx context.Context, request *gen.GetCommentsByPostIDRequest) (*gen.GetCommentsByPostIDResponse, error) {
	slog.Info("GET: GetCommentsByPostID")
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	res := gen.GetCommentsByPostIDResponse{}
	if request.Limit == 0 {
		request.Limit = 10
	}
	comments, err := c.uc.GetCommentsByPostID(ctx, postId, user.ID, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments by post ID")
	}
	for _, comment := range comments {
		res.Comments = append(res.Comments, &gen.CommentHasChildren{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			UserId:          comment.UserID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
			Content:         comment.Content,
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
			Likes: &gen.LikeInfo{
				YourLikedEmoji:    comment.Likes.YourLikedEmoji,
				YourLike:          comment.Likes.YourLike,
				OthersLikedEmojis: comment.Likes.OthersLikedEmojis,
				OthersLikes:       comment.Likes.OthersLikes,
			},
			Attachments: lo.Map(comment.Attachments, func(item *domainUpload.Attachment, _ int) *gen.Attachment {
				return &gen.Attachment{
					Id:             item.ID.String(),
					UserId:         item.UserID.String(),
					AttachableType: item.AttachableType,
					AttachableId:   item.AttachableID.String(),
					Filename:       item.FileName,
					Url:            item.URL,
					UrlThumbnail:   item.URLThumbnail,
					Extension:      item.Extension,
					MimeType:       item.MimeType,
					Folder:         item.Folder,
					CreatedAt:      timestamppb.New(item.CreatedAt),
					UpdatedAt:      timestamppb.New(item.UpdatedAt),
				}
			}),
			Children: lo.Map(comment.Children, func(item *domain.CommentHasMetadata, _ int) *gen.CommentHasMetadata {
				return &gen.CommentHasMetadata{
					Id:      item.ID.String(),
					PostId:  item.PostID.String(),
					UserId:  item.UserID.String(),
					ReplyId: item.ReplyID.UUID.String(),
					TagIds:  utils.ConvertArUUIDToArString(item.TagIDs),
					Content: item.Content,
					Likes: &gen.LikeInfo{
						YourLikedEmoji:    comment.Likes.YourLikedEmoji,
						YourLike:          comment.Likes.YourLike,
						OthersLikedEmojis: comment.Likes.OthersLikedEmojis,
						OthersLikes:       comment.Likes.OthersLikes,
					},
					Attachments: lo.Map(comment.Attachments, func(item *domainUpload.Attachment, _ int) *gen.Attachment {
						return &gen.Attachment{
							Id:             item.ID.String(),
							UserId:         item.UserID.String(),
							AttachableType: item.AttachableType,
							AttachableId:   item.AttachableID.String(),
							Filename:       item.FileName,
							Url:            item.URL,
							UrlThumbnail:   item.URLThumbnail,
							Extension:      item.Extension,
							MimeType:       item.MimeType,
							Folder:         item.Folder,
							CreatedAt:      timestamppb.New(item.CreatedAt),
							UpdatedAt:      timestamppb.New(item.UpdatedAt),
						}
					}),
					ParentCommentId: item.ParentCommentID.UUID.String(),
					CreatedAt:       timestamppb.New(item.CreatedAt),
					UpdatedAt:       timestamppb.New(item.UpdatedAt),
				}
			}),
		})
	}
	return &res, nil
}

// UpdateComment implements gen.CommentServiceServer.
func (c *commentGRPCServer) UpdateComment(ctx context.Context, request *gen.UpdateCommentRequest) (*gen.UpdateCommentResponse, error) {
	slog.Info("PUT: UpdateComment")
	replyToId, _ := uuid.Parse(request.Comment.ReplyId)
	tagIds, _ := utils.ConvertArStringToArUUID(request.Comment.TagIds)
	model := domain.Comment{
		Content: request.Comment.Content,
		ReplyID: uuid.NullUUID{
			UUID:  replyToId,
			Valid: request.Comment.ReplyId != "",
		},
		TagIDs: tagIds,
	}
	comment, err := c.uc.UpdateComment(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.UpdateComment failed")
	}
	res := &gen.UpdateCommentResponse{
		Comment: &gen.Comment{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
			Content:         comment.Content,
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			UserId:          comment.UserID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
		},
	}
	return res, nil
}
