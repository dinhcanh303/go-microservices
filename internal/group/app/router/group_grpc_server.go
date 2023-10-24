package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
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

type groupGRPCServer struct {
	gen.UnimplementedGroupServiceServer
	cfg           *config.Config
	ucGroup       groups.UseCase
	ucGroupMember groupmembers.UseCase
}

var _ gen.GroupServiceServer = (*groupGRPCServer)(nil)

var GroupGRPCServerSet = wire.NewSet(NewGRPCGroupServer)

func NewGRPCGroupServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	ucGroup groups.UseCase,
	ucGroupMember groupmembers.UseCase,
) gen.GroupServiceServer {
	svc := groupGRPCServer{
		cfg:           cfg,
		ucGroup:       ucGroup,
		ucGroupMember: ucGroupMember,
	}
	gen.RegisterGroupServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *groupGRPCServer) GetAllGroupMembers(ctx context.Context, request *gen.GetAllGroupMembersRequest) (*gen.GetAllGroupMembersResponse, error) {
	slog.Info("GET: GetAllGroupMembers")
	groupId, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	groupMembers, err := g.ucGroupMember.GetAllGroupMembers(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroupMember.GetAllGroupMembers failed")
	}
	return &gen.GetAllGroupMembersResponse{
		GroupMembers: lo.Map(groupMembers, func(groupMember *domain.GroupMember, _ int) *gen.GroupMemberResponse {
			return &gen.GroupMemberResponse{
				Id:        groupMember.ID.String(),
				GroupId:   groupMember.GroupID.String(),
				UserId:    groupMember.UserID.String(),
				Role:      groupMember.Role,
				CreatedAt: timestamppb.New(groupMember.CreatedAt),
				UpdatedAt: timestamppb.New(groupMember.UpdatedAt),
			}
		}),
	}, nil
}
func (g *groupGRPCServer) CreateGroupMember(ctx context.Context, request *gen.CreateGroupMemberRequest) (*gen.CreateGroupMemberResponse, error) {
	slog.Info("POST: CreateGroupMember")
	userId, err := uuid.Parse(request.GroupMember.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	groupId, err := uuid.Parse(request.GroupMember.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.GroupMember{
		ID:      uuid.New(),
		GroupID: groupId,
		UserID:  userId,
		Role:    request.GroupMember.Role,
	}
	groupMember, err := g.ucGroupMember.CreateGroupMember(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroupMember.CreateGroupMember failed")
	}
	return &gen.CreateGroupMemberResponse{
		GroupMember: &gen.GroupMemberResponse{
			Id:        groupMember.ID.String(),
			GroupId:   groupMember.GroupID.String(),
			UserId:    groupMember.UserID.String(),
			Role:      groupMember.Role,
			CreatedAt: timestamppb.New(groupMember.CreatedAt),
			UpdatedAt: timestamppb.New(groupMember.UpdatedAt),
		},
	}, nil

}
func (g *groupGRPCServer) DeleteGroupMember(ctx context.Context, request *gen.DeleteGroupMemberRequest) (*gen.DeleteGroupMemberResponse, error) {
	slog.Info("DELETE: DeleteGroupMember")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.ucGroupMember.DeleteGroupMember(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
	}
	return &gen.DeleteGroupMemberResponse{
		Deleted: deleted,
	}, nil
}
func (g *groupGRPCServer) UpdateGroupMember(ctx context.Context, request *gen.UpdateGroupMemberRequest) (*gen.UpdateGroupMemberResponse, error) {
	slog.Info("PUT: UpdateGroupMember")
	id, err := uuid.Parse(request.GroupMember.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.GroupMember{
		ID:   id,
		Role: request.GroupMember.Role,
	}
	groupMember, err := g.ucGroupMember.UpdateGroupMember(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.UpdateGroupMember failed")
	}
	res := &gen.UpdateGroupMemberResponse{
		GroupMember: &gen.GroupMemberResponse{
			Id:        groupMember.ID.String(),
			GroupId:   groupMember.GroupID.String(),
			UserId:    groupMember.UserID.String(),
			Role:      groupMember.Role,
			CreatedAt: timestamppb.New(groupMember.CreatedAt),
			UpdatedAt: timestamppb.New(groupMember.UpdatedAt),
		},
	}
	return res, nil
}

func (g *groupGRPCServer) GetAllGroupByUserId(ctx context.Context, request *gen.GetAllGroupByUserIdRequest) (*gen.GetAllGroupByUserIdResponse, error) {
	slog.Info("GET: GetAllGroupByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groups, err := g.ucGroup.GetAllGroupByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetAllGroupByUserId failed")
	}
	return &gen.GetAllGroupByUserIdResponse{
		Groups: lo.Map(groups, func(group *domain.Group, _ int) *gen.GroupResponse {
			return &gen.GroupResponse{
				Id:          group.ID.String(),
				Name:        group.Name,
				Description: group.Description,
				Status:      group.Status,
				UserId:      group.UserID.String(),
				CreatedAt:   timestamppb.New(group.CreatedAt),
				UpdatedAt:   timestamppb.New(group.UpdatedAt),
			}
		}),
	}, nil
}

func (g *groupGRPCServer) GetAllGroupIdByUserId(ctx context.Context, request *gen.GetAllGroupIdByUserIdRequest) (*gen.GetAllGroupIdByUserIdResponse, error) {
	slog.Info("GET: GetAllGroupByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groupIds, err := g.ucGroup.GetAllGroupIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetAllGroupIdByUserId failed")
	}
	return &gen.GetAllGroupIdByUserIdResponse{
		GroupIds: groupIds,
	}, nil
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
	group, err := g.ucGroup.CreateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.CreateGroup failed")
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
	group, err := g.ucGroup.GetGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
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
	deleted, err := g.ucGroup.DeleteGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
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

	group, err := g.ucGroup.UpdateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.UpdateGroup failed")
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
