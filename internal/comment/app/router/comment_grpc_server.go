package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type commentGRPCServer struct {
	gen.UnimplementedCommentServiceServer
	cfg *config.Config
	uc  posts.UseCase
}

var _ gen.CommentServiceServer = (*commentGRPCServer)(nil)

var CommentGRPCServerSet = wire.NewSet(NewGRPCCommentServer)

func NewGRPCCommentServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc posts.UseCase,
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
func (*commentGRPCServer) CountCommentByCommentID(context.Context, *gen.CountCommentByCommentIDRequest) (*gen.CountCommentByCommentIDResponse, error) {
	panic("unimplemented")
}

// CreateComment implements gen.CommentServiceServer.
func (*commentGRPCServer) CreateComment(context.Context, *gen.CreateCommentRequest) (*gen.CreateCommentResponse, error) {
	panic("unimplemented")
}

// DeleteComment implements gen.CommentServiceServer.
func (*commentGRPCServer) DeleteComment(context.Context, *gen.DeleteCommentRequest) (*gen.DeleteCommentResponse, error) {
	panic("unimplemented")
}

// GetComment implements gen.CommentServiceServer.
func (*commentGRPCServer) GetComment(context.Context, *gen.GetCommentRequest) (*gen.GetCommentResponse, error) {
	panic("unimplemented")
}

// ListCommentByPostID implements gen.CommentServiceServer.
func (*commentGRPCServer) GetCommentsByPostID(context.Context, *gen.GetCommentsByPostIDRequest) (*gen.GetCommentsByPostIDResponse, error) {
	panic("unimplemented")
}

// UpdateComment implements gen.CommentServiceServer.
func (*commentGRPCServer) UpdateComment(context.Context, *gen.UpdateCommentRequest) (*gen.UpdateCommentResponse, error) {
	panic("unimplemented")
}
