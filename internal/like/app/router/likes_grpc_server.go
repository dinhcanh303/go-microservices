package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type likeGRPCServer struct {
	gen.UnimplementedLikeServiceServer
	cfg *config.Config
	uc  likes.UseCase
}

var _ gen.LikeServiceServer = (*likeGRPCServer)(nil)

var LikeGRPCServerSet = wire.NewSet(NewGRPCLikeServer)

func NewGRPCLikeServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc likes.UseCase,
) gen.LikeServiceServer {
	svc := likeGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	gen.RegisterLikeServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (l *likeGRPCServer) GetLikesByPostID(ctx context.Context, request *gen.GetLikesByPostIDRequest) (*gen.GetLikesByPostIDResponse, error) {
	slog.Info("GET: GetLikesByPostID")
	postId, err := uuid.Parse(request.PostID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	likes, err := l.uc.GetLikesByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetLikesByPostID failed")
	}
	res := &gen.GetLikesByPostIDResponse{
		Likes: lo.Map(likes, func(item *domain.Like, _ int) *gen.LikeResponse {
			return &gen.LikeResponse{
				Id:           item.ID.String(),
				UserId:       item.UserID.String(),
				Emoji:        item.Emoji,
				LikeableType: item.LikeableType,
				LikeableId:   item.LikeableID.String(),
				CreatedAt:    timestamppb.New(item.CreatedAt),
				UpdatedAt:    timestamppb.New(item.UpdatedAt),
			}
		}),
	}
	return res, nil
}

func (l *likeGRPCServer) GetLikesByCommentID(ctx context.Context, request *gen.GetLikesByCommentIDRequest) (*gen.GetLikesByCommentIDResponse, error) {
	slog.Info("GET: GetLikesByCommentID")
	commentId, err := uuid.Parse(request.CommentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	likes, err := l.uc.GetLikesByCommentID(ctx, commentId)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetLikesByCommentID failed")
	}
	res := &gen.GetLikesByCommentIDResponse{
		Likes: lo.Map(likes, func(item *domain.Like, _ int) *gen.LikeResponse {
			return &gen.LikeResponse{
				Id:           item.ID.String(),
				UserId:       item.UserID.String(),
				Emoji:        item.Emoji,
				LikeableType: item.LikeableType,
				LikeableId:   item.LikeableID.String(),
				CreatedAt:    timestamppb.New(item.CreatedAt),
				UpdatedAt:    timestamppb.New(item.UpdatedAt),
			}
		}),
	}
	return res, nil
}

// CreateLike implements gen.LikeServiceServer.
func (l *likeGRPCServer) CreateLike(ctx context.Context, request *gen.CreateLikeRequest) (*gen.CreateLikeResponse, error) {
	slog.Info("POST: CreateLike")
	typeLike := []string{"Like/Comment", "Like/Post"}
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
		return nil, errors.Wrap(err, "uc.CreateLike failed")
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
}

// DeleteLike implements gen.LikeServiceServer.
func (l *likeGRPCServer) DeleteLike(ctx context.Context, request *gen.DeleteLikeRequest) (*gen.DeleteLikeResponse, error) {
	slog.Info("POST: UpdateLike")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	like, err := l.uc.DeleteLike(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.DeleteLike failed")
	}
	res := &gen.DeleteLikeResponse{
		Deleted: like,
	}
	return res, nil
}

// UpdateLike implements gen.LikeServiceServer.
func (l *likeGRPCServer) UpdateLike(ctx context.Context, request *gen.UpdateLikeRequest) (*gen.UpdateLikeResponse, error) {
	slog.Info("POST: UpdateLike")
	id, err := uuid.Parse(request.Like.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.Like{
		ID:    id,
		Emoji: request.Like.Emoji,
	}
	like, err := l.uc.UpdateLike(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.UpdateLike failed")
	}
	res := &gen.UpdateLikeResponse{
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
}
