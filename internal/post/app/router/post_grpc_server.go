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

type postGRPCServer struct {
	gen.UnimplementedPostServiceServer
	cfg                  *config.Config
	uc                   posts.UseCase
	uploadDomainService  domain.UploadDomainService
	commentDomainService domain.CommentDomainService
	likeDomainService    domain.LikeDomainService
	groupDomainService   domain.GroupDomainService
	authDomainService    domain.AuthDomainService
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
	groupDomainService domain.GroupDomainService,
	authDomainService domain.AuthDomainService,
) gen.PostServiceServer {
	svc := postGRPCServer{
		cfg:                  cfg,
		uc:                   uc,
		uploadDomainService:  uploadDomainService,
		commentDomainService: commentDomainService,
		likeDomainService:    likeDomainService,
		groupDomainService:   groupDomainService,
		authDomainService:    authDomainService,
	}
	gen.RegisterPostServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (p *postGRPCServer) CreatePost(ctx context.Context, request *gen.CreatePostRequest) (*gen.CreatePostResponse, error) {
	slog.Info("POST: CreatePost")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	groupId, _ := uuid.Parse(request.Post.GroupId)
	model := domain.Post{
		Title:   request.Post.Title,
		Content: request.Post.Content,
		Status:  request.Post.Status,
		UserID:  user.ID,
		GroupID: uuid.NullUUID{
			UUID:  groupId,
			Valid: request.Post.GroupId != "",
		},
	}
	post, err := p.uc.CreatePost(ctx, &model)
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
func (p *postGRPCServer) GetPost(ctx context.Context, request *gen.GetPostRequest) (*gen.GetPostResponse, error) {
	slog.Info("GET: GetPost")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	post, err := p.uc.GetPost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}
	posts := make([]*domain.Post, 0)
	posts = append(posts, post)
	results := manyPostResponse(posts, p, ctx)
	return results[0], nil
}

func (p *postGRPCServer) NewFeed(ctx context.Context, request *gen.NewFeedRequest) (*gen.NewFeedResponse, error) {
	slog.Info("GET: NewFeed")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	userIds, err := p.authDomainService.GetAllUserIdByUserId(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user id service auth")
	}
	slog.Info("U IDS::", userIds)
	groupIds, err := p.groupDomainService.GetAllGroupIdByUserId(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user id service group")
	}
	slog.Info("G IDS::", groupIds)
	posts, err := p.uc.GetPostsByFeed(ctx, userIds, groupIds, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &gen.NewFeedResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) GetPostsByUserId(ctx context.Context, request *gen.GetPostsByUserIdRequest) (*gen.GetPostsByUserIdResponse, error) {
	slog.Info("GET: GetPostsByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse failed")
	}
	posts, err := p.uc.GetPostsByUserId(ctx, userId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByUserId failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &gen.GetPostsByUserIdResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) GetPostsByGroupId(ctx context.Context, request *gen.GetPostsByGroupIdRequest) (*gen.GetPostsByGroupIdResponse, error) {
	slog.Info("GET: GetPostsByGroupId")
	groupId, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse failed")
	}
	posts, err := p.uc.GetPostsByGroupId(ctx, groupId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &gen.GetPostsByGroupIdResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) DeletePost(ctx context.Context, request *gen.DeletePostRequest) (*gen.DeletePostResponse, error) {
	slog.Info("DELETE: DeletePost")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := p.uc.DeletePost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}

	return &gen.DeletePostResponse{
		Deleted: deleted,
	}, nil
}
func (p *postGRPCServer) UpdatePost(ctx context.Context, request *gen.UpdatePostRequest) (*gen.UpdatePostResponse, error) {
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
	post, err := p.uc.UpdatePost(ctx, &model)
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

// private function
func manyPostResponse(posts []*domain.Post, p *postGRPCServer, ctx context.Context) []*gen.GetPostResponse {
	results := make([]*gen.GetPostResponse, 0)
	// channel := make(chan *gen.GetPostResponse, len(posts))
	// var wg sync.WaitGroup
	for _, post := range posts {
		// wg.Add(1)
		likes, err := p.likeDomainService.GetLikesByPostID(ctx, post.ID)
		if err != nil {
			slog.Warn("likeDomainService.GetLikesByPostID failed", err)
			likes = make([]*domainLike.Like, 0)
		}
		comments, err := p.commentDomainService.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			slog.Warn("commentDomainService.GetCommentsByPostID failed", err)
			comments = make([]*sharedkernel.CommentHasChildren, 0)
		}
		attachments, err := p.uploadDomainService.GetAttachmentsByType(ctx, "Attachment/Post", post.ID)
		if err != nil {
			slog.Warn("uploadDomainService.GetAttachmentsByType failed", err)
			attachments = make([]*domainUpload.Attachment, 0)
		}
		results = append(results, &gen.GetPostResponse{
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
		})
	}
	return results
}
