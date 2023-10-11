package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postGRPCServer struct {
	gen.UnimplementedPostServiceServer
	cfg *config.Config
	uc  posts.UseCase
}

var _ gen.PostServiceServer = (*postGRPCServer)(nil)

var PostGRPCServerSet = wire.NewSet(NewGRPCPostServer)

func NewGRPCPostServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc posts.UseCase,
) gen.PostServiceServer {
	svc := postGRPCServer{
		cfg: cfg,
		uc:  uc,
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
			Valid: true,
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
