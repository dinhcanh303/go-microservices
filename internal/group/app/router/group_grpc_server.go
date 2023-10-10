package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type groupGRPCServer struct {
	gen.UnimplementedGroupServiceServer
	cfg *config.Config
	uc  groups.UseCase
}

var _ gen.GroupServiceServer = (*groupGRPCServer)(nil)

var GroupGRPCServerSet = wire.NewSet(NewGRPCGroupServer)

func NewGRPCGroupServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc groups.UseCase,
) gen.GroupServiceServer {
	svc := groupGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	gen.RegisterGroupServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *groupGRPCServer) CreateGroup(ctx context.Context, request *gen.CreateGroupRequest) (*gen.CreateGroupResponse, error) {
	slog.Info("POST: CreateGroup")
	userId, err := uuid.Parse(request.Group.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.Group{
		Name:        request.Group.Name,
		Description: request.Group.Description,
		Status:      request.Group.Status,
		UserID:      userId,
	}
	slog.Info("Model", model)

	group, err := g.uc.CreateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreateGroup failed")
	}
	res := &gen.CreateGroupResponse{
		Group: &gen.GroupResponse{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
	}
	return res, nil
}
func (g *groupGRPCServer) GetGroup(ctx context.Context, request *gen.GetGroupRequest) (*gen.GetGroupResponse, error) {
	slog.Info("GET: GetGroup")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	group, err := g.uc.GetGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetGroup failed")
	}
	res := &gen.GetGroupResponse{
		Group: &gen.GroupResponse{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
	}
	return res, nil
}
func (g *groupGRPCServer) DeleteGroup(ctx context.Context, request *gen.DeleteGroupRequest) (*gen.DeleteGroupResponse, error) {
	slog.Info("DELETE: DeleteGroup")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.uc.DeleteGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetGroup failed")
	}

	return &gen.DeleteGroupResponse{
		Deleted: deleted,
	}, nil
}
func (g *groupGRPCServer) UpdateGroup(ctx context.Context, request *gen.UpdateGroupRequest) (*gen.UpdateGroupResponse, error) {
	slog.Info("PUT: UpdateGroup")
	id, err := uuid.Parse(request.Group.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.Group{
		ID:          id,
		Name:        request.Group.Name,
		Description: request.Group.Description,
		Status:      request.Group.Status,
	}

	group, err := g.uc.UpdateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreateGroup failed")
	}
	res := &gen.UpdateGroupResponse{
		Group: &gen.GroupResponse{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
	}
	return res, nil
}
