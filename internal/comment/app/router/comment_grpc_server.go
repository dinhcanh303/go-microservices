package router

import (
	"context"

	v1a "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1 "github.com/dinhcanh303/go-microservices/api/comment/v1"
	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
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
	v1.UnimplementedCommentServiceServer
	cfg           *config.Config
	uc            comments.UseCase
	authDomainSvc domain.AuthDomainService
}

var _ v1.CommentServiceServer = (*commentGRPCServer)(nil)

var CommentGRPCServerSet = wire.NewSet(NewGRPCCommentServer)

func NewGRPCCommentServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc comments.UseCase,
	authDomainSvc domain.AuthDomainService,
) v1.CommentServiceServer {
	svc := commentGRPCServer{
		cfg:           cfg,
		uc:            uc,
		authDomainSvc: authDomainSvc,
	}
	v1.RegisterCommentServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

// CountCommentByCommentID implements v1.CommentServiceServer.
func (c *commentGRPCServer) CountCommentByCommentID(ctx context.Context, request *v1.CountCommentByCommentIDRequest) (*v1.CountCommentByCommentIDResponse, error) {
	slog.Info("POST: CountCommentByCommentID")
	commentId, err := uuid.Parse(request.CommentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	count, err := c.uc.CountCommentByCommentID(ctx, commentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count")
	}
	return &v1.CountCommentByCommentIDResponse{
		Count: count,
	}, nil
}

// CountCommentByPostID implements v1.CommentServiceServer.
func (c *commentGRPCServer) CountCommentByPostID(ctx context.Context, request *v1.CountCommentByPostIDRequest) (*v1.CountCommentByPostIDResponse, error) {
	slog.Info("POST: CountCommentByPostID")
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	count, err := c.uc.CountCommentByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count")
	}
	return &v1.CountCommentByPostIDResponse{
		Count: count,
	}, nil
}

// CreateComment implements v1.CommentServiceServer.
func (c *commentGRPCServer) CreateComment(ctx context.Context, request *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {
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
	return &v1.CreateCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}

// DeleteComment implements v1.CommentServiceServer.
func (c *commentGRPCServer) DeleteComment(ctx context.Context, request *v1.DeleteCommentRequest) (*v1.DeleteCommentResponse, error) {
	slog.Info("DELETE: DeleteComment")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	result, err := c.uc.DeleteComment(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete")
	}
	return &v1.DeleteCommentResponse{
		Deleted: result,
	}, nil
}

// GetComment implements v1.CommentServiceServer.
func (c *commentGRPCServer) GetComment(ctx context.Context, request *v1.GetCommentRequest) (*v1.GetCommentResponse, error) {
	slog.Info("GET: GetComment")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	comment, err := c.uc.GetComment(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comment")
	}
	return &v1.GetCommentResponse{
		Comment: entityToProtobuf(comment),
	}, nil
}
func (c *commentGRPCServer) GetCommentsByCommentID(ctx context.Context, request *v1.GetCommentsByCommentIDRequest) (*v1.GetCommentsByCommentIDResponse, error) {
	slog.Info("GET: GetCommentsByCommentID")
	commentId, err := uuid.Parse(request.CommentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	res := v1.GetCommentsByCommentIDResponse{}
	if request.Limit == 0 {
		request.Limit = 10
	}
	comments, err := c.uc.GetCommentsByCommentID(ctx, commentId, user.ID, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments by post ID")
	}
	for _, comment := range comments {
		user, err := c.authDomainSvc.GetProfile(ctx, comment.UserID)
		if err != nil {
			user = &v1a.GetProfileResponse{}
		}
		replyName := ""
		if comment.UserID == comment.ReplyID.UUID {
			replyName = user.User.FullName
		} else {
			if comment.ReplyID.UUID.String() != constant.NullUUID {
				replyUser, err := c.authDomainSvc.GetProfile(ctx, comment.ReplyID.UUID)
				if err == nil {
					replyName = replyUser.User.FullName
				}
			}
		}
		tagIds, tagNames := handleTags(ctx, comment.TagIDs, c.authDomainSvc)
		res.Comments = append(res.Comments, entityCommentToProtobuf(comment, replyName, tagIds, tagNames, user.User))
	}
	return &res, nil
}

// ListCommentByPostID implements v1.CommentServiceServer.
func (c *commentGRPCServer) GetCommentsByPostID(ctx context.Context, request *v1.GetCommentsByPostIDRequest) (*v1.GetCommentsByPostIDResponse, error) {
	slog.Info("GET: GetCommentsByPostID")
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	res := v1.GetCommentsByPostIDResponse{}
	if request.Limit == 0 {
		request.Limit = 10
	}
	comments, err := c.uc.GetCommentsByPostID(ctx, postId, user.ID, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments by post ID")
	}
	for _, comment := range comments {
		user, err := c.authDomainSvc.GetProfile(ctx, comment.UserID)
		if err != nil {
			user = &v1a.GetProfileResponse{}
		}
		replyName := ""
		if comment.UserID == comment.ReplyID.UUID {
			replyName = user.User.FullName
		} else {
			if comment.ReplyID.UUID.String() != constant.NullUUID {
				replyUser, err := c.authDomainSvc.GetProfile(ctx, comment.ReplyID.UUID)
				if err == nil {
					replyName = replyUser.User.FullName
				}
			}
		}
		tagIds, tagNames := handleTags(ctx, comment.TagIDs, c.authDomainSvc)
		res.Comments = append(res.Comments, &v1.CommentHasChildren{
			Id:              comment.ID.String(),
			PostId:          comment.PostID.String(),
			UserId:          comment.UserID.String(),
			ReplyId:         comment.ReplyID.UUID.String(),
			ReplyName:       replyName,
			TagIds:          tagIds,
			TagNames:        tagNames,
			Content:         comment.Content,
			ParentCommentId: comment.ParentCommentID.UUID.String(),
			CreatedAt:       timestamppb.New(comment.CreatedAt),
			UpdatedAt:       timestamppb.New(comment.UpdatedAt),
			User:            user.User,
			Likes:           sharedkernel.EntityLikeToProtobuf(comment.Likes),
			Attachments:     sharedkernel.EntityAttachmentToProtobuf(comment.Attachments),
			Children: lo.Map(comment.Children, func(item *domain.CommentHasMetadata, _ int) *v1.CommentHasMetadata {
				user, err := c.authDomainSvc.GetProfile(ctx, item.UserID)
				if err != nil {
					user = &v1a.GetProfileResponse{}
				}
				replyName := ""
				if item.UserID == item.ReplyID.UUID {
					replyName = user.User.FullName
				} else {
					if item.ReplyID.UUID.String() != constant.NullUUID {
						replyUser, err := c.authDomainSvc.GetProfile(ctx, item.ReplyID.UUID)
						if err == nil {
							replyName = replyUser.User.FullName
						}
					}
				}
				tagIds, tagNames := handleTags(ctx, item.TagIDs, c.authDomainSvc)
				return entityCommentToProtobuf(item, replyName, tagIds, tagNames, user.User)
			}),
		})
	}
	return &res, nil
}

// UpdateComment implements v1.CommentServiceServer.
func (c *commentGRPCServer) UpdateComment(ctx context.Context, request *v1.UpdateCommentRequest) (*v1.UpdateCommentResponse, error) {
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
	res := &v1.UpdateCommentResponse{
		Comment: entityToProtobuf(comment),
	}
	return res, nil
}
func handleTags(ctx context.Context, tagIds []uuid.UUID, authDomainSvc domain.AuthDomainService) ([]string, []string) {
	strTagIds := utils.ConvertArUUIDToArString(tagIds)
	tagNames := make([]string, 0)
	for _, tagId := range tagIds {
		user, _ := authDomainSvc.GetProfile(ctx, tagId)
		tagNames = append(tagNames, user.User.FullName)
	}
	return strTagIds, tagNames
}
func entityToProtobuf(comment *domain.Comment) *v1.Comment {
	return &v1.Comment{
		Id:              comment.ID.String(),
		PostId:          comment.PostID.String(),
		ReplyId:         comment.ReplyID.UUID.String(),
		TagIds:          utils.ConvertArUUIDToArString(comment.TagIDs),
		Content:         comment.Content,
		ParentCommentId: comment.ParentCommentID.UUID.String(),
		UserId:          comment.UserID.String(),
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}

func entityCommentToProtobuf(comment *domain.CommentHasMetadata, replyName string, tagIds []string, tagNames []string, user *v1a.User) *v1.CommentHasMetadata {
	return &v1.CommentHasMetadata{
		Id:              comment.ID.String(),
		PostId:          comment.PostID.String(),
		UserId:          comment.UserID.String(),
		ReplyId:         comment.ReplyID.UUID.String(),
		ReplyName:       replyName,
		TagIds:          tagIds,
		TagNames:        tagNames,
		Content:         comment.Content,
		User:            user,
		Likes:           sharedkernel.EntityLikeToProtobuf(comment.Likes),
		Attachments:     sharedkernel.EntityAttachmentToProtobuf(comment.Attachments),
		ParentCommentId: comment.ParentCommentID.UUID.String(),
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}
