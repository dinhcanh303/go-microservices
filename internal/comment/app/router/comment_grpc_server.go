package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/comment/config"
	"github.com/dinhcanh303/go-microservices/internal/comment/usecases/comments"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
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
func (*commentGRPCServer) CountCommentByCommentID(context.Context, *gen.CountCommentByCommentIDRequest) (*gen.CountCommentByCommentIDResponse, error) {
	panic("unimplemented")
}

// CreateComment implements gen.CommentServiceServer.
func (*commentGRPCServer) CreateComment(ctx context.Context, request *gen.CreateCommentRequest) (*gen.CreateCommentResponse, error) {
	slog.Info("POST: CreateComment")
	userId, err := uuid.Parse(request.Like.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	likeableId, err := uuid.Parse(request.Like.LikeableId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	likeableType := request.Like.LikeableType
	hasTypeLike := slices.Contains(typeLike, likeableType)
	if !hasTypeLike {
		return nil, errors.Wrap(err, "Please Enter Input Type Correct")
	}
	model := domain.Like{
		ID:           uuid.New(),
		Emoji:        request.Like.Emoji,
		LikeableType: request.Like.LikeableType,
		LikeableID:   likeableId,
		UserID:       userId,
	}
	like, err := l.uc.CreateLike(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreateGroup failed")
	}
	res := &gen.CreateLikeResponse{
		Like: &gen.LikeResponse{
			Id:           like.ID.String(),
			Emoji:        like.Emoji,
			LikeableType: like.LikeableType,
			LikeableId:   like.LikeableID.String(),
			UserId:       like.UserID.String(),
			CreatedAt:    timestamppb.New(like.CreatedAt),
			UpdatedAt:    timestamppb.New(like.UpdatedAt),
		},
	}
	return res, nil
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
