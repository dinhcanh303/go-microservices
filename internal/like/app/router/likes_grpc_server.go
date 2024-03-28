package router

import (
	"context"

	v1 "github.com/dinhcanh303/go-microservices/api/like/v1"
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type likeGRPCServer struct {
	v1.UnimplementedLikeServiceServer
	cfg *config.Config
	uc  likes.UseCase
}

var _ v1.LikeServiceServer = (*likeGRPCServer)(nil)

var LikeGRPCServerSet = wire.NewSet(NewGRPCLikeServer)

func NewGRPCLikeServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc likes.UseCase,
) v1.LikeServiceServer {
	svc := likeGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	v1.RegisterLikeServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (l *likeGRPCServer) GetLikesInfoByPostID(ctx context.Context, request *v1.GetLikesInfoByPostIDRequest) (*v1.GetLikesInfoByPostIDResponse, error) {
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	likeInfo, err := l.uc.GetLikesInfoByPostID(ctx, postId, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetLikesInfoByPostID failed")
	}
	return &v1.GetLikesInfoByPostIDResponse{
		Likes: sharedkernel.EntityLikeToProtobuf(likeInfo),
	}, nil
}

func (l *likeGRPCServer) GetLikesInfoByCommentID(ctx context.Context, request *v1.GetLikesInfoByCommentIDRequest) (*v1.GetLikesInfoByCommentIDResponse, error) {
	commentId, err := uuid.Parse(request.CommentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	likeInfo, err := l.uc.GetLikesInfoByCommentID(ctx, commentId, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetLikesInfoByCommentID failed")
	}
	return &v1.GetLikesInfoByCommentIDResponse{
		Likes: sharedkernel.EntityLikeToProtobuf(likeInfo),
	}, nil
}

// func (l *likeGRPCServer) GetLikesByCommentID(ctx context.Context, request *v1.GetLikesByCommentIDRequest) (*v1.GetLikesByCommentIDResponse, error) {
// 	slog.Info("GET: GetLikesByCommentID")
// 	commentId, err := uuid.Parse(request.CommentId)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to parse")
// 	}
// 	likes, err := l.uc.GetLikesByCommentID(ctx, commentId)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "uc.GetLikesByCommentID failed")
// 	}
// 	res := &v1.GetLikesByCommentIDResponse{
// 		Likes: lo.Map(likes, func(item *domain.Like, _ int) *v1.Like {
// 			return &v1.Like{
// 				Id:           item.ID.String(),
// 				UserId:       item.UserID.String(),
// 				Emoji:        item.Emoji,
// 				LikeableType: item.LikeableType,
// 				LikeableId:   item.LikeableID.String(),
// 				CreatedAt:    timestamppb.New(item.CreatedAt),
// 				UpdatedAt:    timestamppb.New(item.UpdatedAt),
// 			}
// 		}),
// 	}
// 	return res, nil
// }

// CreateLike implements v1.LikeServiceServer.
func (l *likeGRPCServer) CreateLike(ctx context.Context, request *v1.CreateLikeRequest) (*v1.CreateLikeResponse, error) {
	typeLike := []string{constant.LikeCommentType, constant.LikePostType}
	emoji := request.Like.Emoji
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	userId := user.ID
	likeableId, err := uuid.Parse(request.Like.LikeableId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	likeableType := request.Like.LikeableType
	hasTypeLike := slices.Contains(typeLike, likeableType)
	if !hasTypeLike {
		return nil, errors.Wrap(err, "Please Enter Input Type Correct")
	}
	like, _ := l.uc.GetLikeByUserId(ctx, likeableType, likeableId, userId)
	if like != nil {
		if like.Emoji == emoji {
			_, err := l.uc.DeleteLike(ctx, like.ID)
			if err != nil {
				return nil, err
			}
			return &v1.CreateLikeResponse{}, nil
		} else {
			like, err := l.uc.UpdateLike(ctx, &domain.Like{
				ID:    like.ID,
				Emoji: request.Like.Emoji,
			})
			if err != nil {
				return nil, errors.Wrap(err, "uc.UpdateLike failed")
			}
			return &v1.CreateLikeResponse{
				Like: entityToProtobuf(like),
			}, nil
		}
	}
	model := domain.Like{
		ID:           uuid.New(),
		Emoji:        emoji,
		LikeableType: likeableType,
		LikeableID:   likeableId,
		UserID:       userId,
	}
	like, err = l.uc.CreateLike(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreateLike failed")
	}
	res := &v1.CreateLikeResponse{
		Like: entityToProtobuf(like),
	}
	return res, nil
}

// DeleteLike implements v1.LikeServiceServer.
func (l *likeGRPCServer) DeleteLike(ctx context.Context, request *v1.DeleteLikeRequest) (*v1.DeleteLikeResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	like, err := l.uc.DeleteLike(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.DeleteLike failed")
	}
	res := &v1.DeleteLikeResponse{
		Deleted: like,
	}
	return res, nil
}

// UpdateLike implements v1.LikeServiceServer.
func (l *likeGRPCServer) UpdateLike(ctx context.Context, request *v1.UpdateLikeRequest) (*v1.UpdateLikeResponse, error) {
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
	res := &v1.UpdateLikeResponse{
		Like: entityToProtobuf(like),
	}
	return res, nil
}

func entityToProtobuf(like *domain.Like) *v1.Like {
	return &v1.Like{
		Id:           like.ID.String(),
		Emoji:        like.Emoji,
		LikeableType: like.LikeableType,
		LikeableId:   like.LikeableID.String(),
		UserId:       like.UserID.String(),
		CreatedAt:    timestamppb.New(like.CreatedAt),
		UpdatedAt:    timestamppb.New(like.UpdatedAt),
	}
}
