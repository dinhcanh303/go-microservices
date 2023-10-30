package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainComment "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
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

type postGRPCServer struct {
	gen.UnimplementedPostServiceServer
	cfg                  *config.Config
	uc                   posts.UseCase
	uploadDomainService  domain.UploadDomainService
	commentDomainService domain.CommentDomainService
	likeDomainService    domain.LikeDomainService
}

var _ gen.PostServiceServer = (*postGRPCServer)(nil)

var PostGRPCServerSet = wire.NewSet(NewGRPCPostServer)

func NewGRPCPostServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc posts.UseCase,
	uploadDomainService domain.UploadDomainService,
	commentDomainService domain.CommentDomainService,
	likeDomainService domain.LikeDomainService,
) gen.PostServiceServer {
	svc := postGRPCServer{
		cfg:                  cfg,
		uc:                   uc,
		uploadDomainService:  uploadDomainService,
		commentDomainService: commentDomainService,
		likeDomainService:    likeDomainService,
	}
	gen.RegisterPostServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *postGRPCServer) CreatePost(ctx context.Context, request *gen.CreatePostRequest) (*gen.CreatePostResponse, error) {
	slog.Info("POST: CreatePost")
	userId, err := uuid.Parse(request.Post.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "Parse User ID failed")
	}
	groupId, _ := uuid.Parse(request.Post.GroupId)
	slog.Info("GROUP ID:", groupId)
	model := domain.Post{
		Title:   request.Post.Title,
		Content: request.Post.Content,
		Status:  request.Post.Status,
		UserID:  userId,
		GroupID: uuid.NullUUID{
			UUID:  groupId,
			Valid: request.Post.GroupId != "",
		},
	}
	slog.Info("Model", model)

	post, err := g.uc.CreatePost(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreatePost failed")
	}
	res := &gen.CreatePostResponse{
		Post: &gen.PostResponse{
			Id:        post.ID.String(),
			Title:     post.Title,
			Content:   post.Content,
			UserId:    post.UserID.String(),
			GroupId:   post.GroupID.UUID.String(),
			Status:    post.Status,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}
	return res, nil
}
func (g *postGRPCServer) GetPost(ctx context.Context, request *gen.GetPostRequest) (*gen.GetPostResponse, error) {
	slog.Info("GET: GetPost")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	post, err := g.uc.GetPost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}
	likes, err := g.likeDomainService.GetLikesByPostID(ctx, post.ID)
	if err != nil {
		return nil, errors.Wrap(err, "likeDomainService.GetLikesByPostID failed")
	}
	comments, err := g.commentDomainService.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		return nil, errors.Wrap(err, "commentDomainService.GetCommentsByPostID failed")
	}
	attachments, err := g.uploadDomainService.GetAttachmentsByType(ctx, "Attachment/Post", post.ID)
	if err != nil {
		return nil, errors.Wrap(err, "uploadDomainService.GetAttachmentsByType failed")
	}
	res := &gen.GetPostResponse{
		Post: &gen.PostResponse{
			Id:        post.ID.String(),
			Title:     post.Title,
			Content:   post.Content,
			UserId:    post.UserID.String(),
			GroupId:   post.GroupID.UUID.String(),
			Status:    post.Status,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
		Attachments: lo.Map(attachments, func(item *domainUpload.Attachment, _ int) *gen.Attachment {
			return &gen.Attachment{
				Id:             item.ID.String(),
				AttachableType: item.AttachableType,
				AttachableId:   item.AttachableID.String(),
				Filename:       item.FileName,
				Extension:      item.Extension,
				MimeType:       item.MimeType,
				Folder:         item.Folder,
				Url:            item.URL,
				UrlThumbnail:   item.URLThumbnail,
				UserId:         item.UserID.String(),
				CreatedAt:      timestamppb.New(item.CreatedAt),
				UpdatedAt:      timestamppb.New(item.UpdatedAt),
			}
		}),
		Likes: lo.Map(likes, func(item *domainLike.Like, _ int) *gen.Like {
			return &gen.Like{
				Id:           item.ID.String(),
				Emoji:        item.Emoji,
				LikeableType: item.LikeableType,
				LikeableId:   item.LikeableID.String(),
				UserId:       item.UserID.String(),
				CreatedAt:    timestamppb.New(item.CreatedAt),
				UpdatedAt:    timestamppb.New(item.UpdatedAt),
			}
		}),
		Comments: lo.Map(comments, func(item *sharedkernel.CommentHasChildren, _ int) *gen.CommentHasCommentAndLike {
			return &gen.CommentHasCommentAndLike{
				Id:              item.ID.String(),
				PostId:          item.PostID.String(),
				ReplyToId:       item.ReplyToID.UUID.String(),
				Content:         item.Content,
				ParentCommentId: item.ParentCommentID.UUID.String(),
				Likes: lo.Map(item.Likes, func(item *domainLike.Like, _ int) *gen.Like {
					return &gen.Like{
						Id:           item.ID.String(),
						Emoji:        item.Emoji,
						LikeableType: item.LikeableType,
						LikeableId:   item.LikeableID.String(),
						UserId:       item.UserID.String(),
						CreatedAt:    timestamppb.New(item.CreatedAt),
						UpdatedAt:    timestamppb.New(item.UpdatedAt),
					}
				}),
				Children: lo.Map(item.Children, func(item *domainComment.CommentHasLike, _ int) *gen.CommentHasLike {
					return &gen.CommentHasLike{
						Id:              item.ID.String(),
						PostId:          item.PostID.String(),
						UserId:          item.UserID.String(),
						ReplyToId:       item.ID.String(),
						Content:         item.Content,
						ParentCommentId: item.ParentCommentID.UUID.String(),
						Likes: lo.Map(item.Likes, func(like *domainLike.Like, _ int) *gen.Like {
							return &gen.Like{
								Id:           like.ID.String(),
								Emoji:        like.Emoji,
								LikeableType: like.LikeableType,
								LikeableId:   like.LikeableID.String(),
								UserId:       like.UserID.String(),
								CreatedAt:    timestamppb.New(like.CreatedAt),
								UpdatedAt:    timestamppb.New(like.UpdatedAt),
							}
						}),
						CreatedAt: timestamppb.New(item.CreatedAt),
						UpdatedAt: timestamppb.New(item.UpdatedAt),
					}
				}),
				UserId:    item.UserID.String(),
				CreatedAt: timestamppb.New(item.CreatedAt),
				UpdatedAt: timestamppb.New(item.UpdatedAt),
			}
		}),
	}
	return res, nil
}
func (g *postGRPCServer) DeletePost(ctx context.Context, request *gen.DeletePostRequest) (*gen.DeletePostResponse, error) {
	slog.Info("DELETE: DeletePost")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.uc.DeletePost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}

	return &gen.DeletePostResponse{
		Deleted: deleted,
	}, nil
}
func (g *postGRPCServer) UpdatePost(ctx context.Context, request *gen.UpdatePostRequest) (*gen.UpdatePostResponse, error) {
	slog.Info("PUT: UpdatePost")
	id, err := uuid.Parse(request.Post.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.Post{
		ID:      id,
		Title:   request.Post.Title,
		Content: request.Post.Content,
		Status:  request.Post.Status,
	}

	post, err := g.uc.UpdatePost(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreatePost failed")
	}
	res := &gen.UpdatePostResponse{
		Post: &gen.PostResponse{
			Id:        post.ID.String(),
			Title:     post.Title,
			Content:   post.Content,
			UserId:    post.UserID.String(),
			GroupId:   post.GroupID.UUID.String(),
			Status:    post.Status,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}
	return res, nil
}
