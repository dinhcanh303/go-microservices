package router

import (
	"context"
	"sync"

	v1a "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1g "github.com/dinhcanh303/go-microservices/api/group/v1"
	v1 "github.com/dinhcanh303/go-microservices/api/post/v1"
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postGRPCServer struct {
	v1.UnimplementedPostServiceServer
	cfg                  *config.Config
	uc                   posts.UseCase
	uploadDomainService  domain.UploadDomainService
	commentDomainService domain.CommentDomainService
	likeDomainService    domain.LikeDomainService
	groupDomainService   domain.GroupDomainService
	authDomainService    domain.AuthDomainService
}

var _ v1.PostServiceServer = (*postGRPCServer)(nil)

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
) v1.PostServiceServer {
	svc := postGRPCServer{
		cfg:                  cfg,
		uc:                   uc,
		uploadDomainService:  uploadDomainService,
		commentDomainService: commentDomainService,
		likeDomainService:    likeDomainService,
		groupDomainService:   groupDomainService,
		authDomainService:    authDomainService,
	}
	v1.RegisterPostServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (p *postGRPCServer) CreatePost(ctx context.Context, request *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	groupId, _ := uuid.Parse(request.Post.GroupId)
	model := domain.Post{
		Content:   request.Post.Content,
		BgContent: request.Post.BgContent,
		Status:    request.Post.Status,
		UserID:    user.ID,
		GroupID: uuid.NullUUID{
			UUID:  groupId,
			Valid: request.Post.GroupId != "",
		},
	}
	post, err := p.uc.CreatePost(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreatePost failed")
	}
	res := &v1.CreatePostResponse{
		Post: entityToProtobuf(post),
	}
	return res, nil
}
func (p *postGRPCServer) GetPost(ctx context.Context, request *v1.GetPostRequest) (*v1.GetPostResponse, error) {
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
func (p *postGRPCServer) GetPostNormal(ctx context.Context, request *v1.GetPostNormalRequest) (*v1.GetPostNormalResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	post, err := p.uc.GetPost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}
	return &v1.GetPostNormalResponse{
		Post: entityToProtobuf(post),
	}, nil
}

func (p *postGRPCServer) NewFeed(ctx context.Context, request *v1.NewFeedRequest) (*v1.NewFeedResponse, error) {
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	userIds, err := p.authDomainService.GetUserIdsByUserId(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user id service auth")
	}
	groupIds, err := p.groupDomainService.GetGroupIdsByUserId(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user id service group")
	}
	posts, err := p.uc.GetPostsByFeed(ctx, userIds, groupIds, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &v1.NewFeedResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) NewFeedGroups(ctx context.Context, request *v1.NewFeedGroupsRequest) (*v1.NewFeedGroupsResponse, error) {
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	groupIds, err := p.groupDomainService.GetGroupIdsByUserId(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get user id service group")
	}
	posts, err := p.uc.GetPostsByFeedGroup(ctx, groupIds, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &v1.NewFeedGroupsResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) GetPostsByUserId(ctx context.Context, request *v1.GetPostsByUserIdRequest) (*v1.GetPostsByUserIdResponse, error) {
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse failed")
	}
	posts, err := p.uc.GetPostsByUserId(ctx, userId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByUserId failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &v1.GetPostsByUserIdResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) GetPostsByGroupId(ctx context.Context, request *v1.GetPostsByGroupIdRequest) (*v1.GetPostsByGroupIdResponse, error) {
	groupId, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse failed")
	}
	posts, err := p.uc.GetPostsByGroupId(ctx, groupId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
	}
	results := manyPostResponse(posts, p, ctx)
	return &v1.GetPostsByGroupIdResponse{
		Posts: results,
	}, nil
}

func (p *postGRPCServer) DeletePost(ctx context.Context, request *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := p.uc.DeletePost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPost failed")
	}

	return &v1.DeletePostResponse{
		Deleted: deleted,
	}, nil
}
func (p *postGRPCServer) UpdatePost(ctx context.Context, request *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	id, err := uuid.Parse(request.Post.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.Post{
		ID:        id,
		Content:   request.Post.Content,
		BgContent: request.Post.BgContent,
		Status:    request.Post.Status,
	}
	post, err := p.uc.UpdatePost(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreatePost failed")
	}
	res := &v1.UpdatePostResponse{
		Post: entityToProtobuf(post),
	}
	return res, nil
}

// private function
func manyPostResponse(posts []*domain.Post, p *postGRPCServer, ctx context.Context) []*v1.GetPostResponse {
	results := make([]*v1.GetPostResponse, len(posts))
	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	for i, post := range posts {
		wg.Add(1)
		go func(index int, post *domain.Post) {
			defer wg.Done()
			likeInfoCh := make(chan *domainLike.LikesInfo, 1)
			countCommentsCh := make(chan int64, 1)
			userCh := make(chan *v1a.GetProfileResponse, 1)
			groupCh := make(chan *v1g.GetGroupResponse, 1)
			attachmentsCh := make(chan []*domainUpload.Attachment, 1)
			go func() {
				likeInfo, err := p.likeDomainService.GetLikesByPostID(ctx, post.ID)
				if err != nil {
					slog.Warn("likeDomainService.GetLikesByPostID failed", err)
					likeInfo = &domainLike.LikesInfo{}
				}
				likeInfoCh <- likeInfo
			}()
			go func() {
				countComments, err := p.commentDomainService.CountCommentByPostID(ctx, post.ID)
				if err != nil {
					slog.Warn("commentDomainService.CountCommentByPostID failed", err)
					countComments = 0
				}
				countCommentsCh <- countComments
			}()
			go func() {
				user, err := p.authDomainService.GetProfile(ctx, post.UserID)
				if err != nil {
					user = &v1a.GetProfileResponse{}
				}
				userCh <- user
			}()
			go func() {
				group, err := p.groupDomainService.GetGroup(ctx, post.GroupID)
				if err != nil {
					group = &v1g.GetGroupResponse{}
				}
				groupCh <- group
			}()
			go func() {
				attachments, err := p.uploadDomainService.GetAttachmentsByType(ctx, constant.ATTACHMENT_POST, post.ID)
				if err != nil {
					slog.Warn("uploadDomainService.GetAttachmentsByType failed", err)
					attachments = make([]*domainUpload.Attachment, 0)
				}
				attachmentsCh <- attachments
			}()
			likeInfo := <-likeInfoCh
			countComments := <-countCommentsCh
			user := <-userCh
			group := <-groupCh
			attachments := <-attachmentsCh
			result := &v1.GetPostResponse{
				Post: &v1.Post{
					Id:        post.ID.String(),
					Content:   post.Content,
					BgContent: post.BgContent,
					UserId:    post.UserID.String(),
					GroupId:   post.GroupID.UUID.String(),
					Status:    post.Status,
					CreatedAt: timestamppb.New(post.CreatedAt),
					UpdatedAt: timestamppb.New(post.UpdatedAt),
				},
				Group:         group.Group,
				User:          user.User,
				CountComments: countComments,
				Attachments:   sharedkernel.EntityAttachmentToProtobuf(attachments),
				Likes:         sharedkernel.EntityLikeToProtobuf(likeInfo),
			}
			mutex.Lock()
			results[index] = result
			mutex.Unlock()
		}(i, post)
	}
	wg.Wait()
	return results
}

func entityToProtobuf(post *domain.Post) *v1.Post {
	return &v1.Post{
		Id:        post.ID.String(),
		Content:   post.Content,
		BgContent: post.BgContent,
		UserId:    post.UserID.String(),
		GroupId:   post.GroupID.UUID.String(),
		Status:    post.Status,
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
	}
}
