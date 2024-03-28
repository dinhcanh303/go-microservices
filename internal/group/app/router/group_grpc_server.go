package router

import (
	"context"

	v1a "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1 "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
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
	v1.UnimplementedGroupServiceServer
	cfg               *config.Config
	ucGroup           groups.UseCase
	ucGroupMember     groupmembers.UseCase
	authDomainService domain.AuthDomainService
}

var _ v1.GroupServiceServer = (*groupGRPCServer)(nil)

var GroupGRPCServerSet = wire.NewSet(NewGRPCGroupServer)

func NewGRPCGroupServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	ucGroup groups.UseCase,
	ucGroupMember groupmembers.UseCase,
	authDomainService domain.AuthDomainService,
) v1.GroupServiceServer {
	svc := groupGRPCServer{
		cfg:               cfg,
		ucGroup:           ucGroup,
		ucGroupMember:     ucGroupMember,
		authDomainService: authDomainService,
	}
	v1.RegisterGroupServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *groupGRPCServer) GetGroupMembers(ctx context.Context, request *v1.GetGroupMembersRequest) (*v1.GetGroupMembersResponse, error) {
	slog.Info("GET: GetAllGroupMembers")
	groupId, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	groupMembers, err := g.ucGroupMember.GetGroupMembers(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroupMember.GetAllGroupMembers failed")
	}
	return &v1.GetGroupMembersResponse{
		GroupMembers: lo.Map(groupMembers, func(groupMember *domain.GroupMember, _ int) *v1.GroupMemberMetadata {
			user, err := g.authDomainService.GetProfile(ctx, groupMember.UserID)
			if err != nil {
				user = &v1a.GetProfileResponse{}
			}
			return &v1.GroupMemberMetadata{
				Id:        groupMember.ID.String(),
				GroupId:   groupMember.GroupID.String(),
				UserId:    groupMember.UserID.String(),
				User:      user.User,
				Role:      groupMember.Role,
				CreatedAt: timestamppb.New(groupMember.CreatedAt),
				UpdatedAt: timestamppb.New(groupMember.UpdatedAt),
			}
		}),
	}, nil
}
func (g *groupGRPCServer) CreateGroupMember(ctx context.Context, request *v1.CreateGroupMemberRequest) (*v1.CreateGroupMemberResponse, error) {
	slog.Info("POST: CreateGroupMember")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	groupId, err := uuid.Parse(request.GroupMember.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	model := domain.GroupMember{
		ID:      uuid.New(),
		GroupID: groupId,
		UserID:  user.ID,
		Role:    request.GroupMember.Role,
	}
	groupMember, err := g.ucGroupMember.CreateGroupMember(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroupMember.CreateGroupMember failed")
	}
	return &v1.CreateGroupMemberResponse{
		GroupMember: &v1.GroupMember{
			Id:        groupMember.ID.String(),
			GroupId:   groupMember.GroupID.String(),
			UserId:    groupMember.UserID.String(),
			Role:      groupMember.Role,
			CreatedAt: timestamppb.New(groupMember.CreatedAt),
			UpdatedAt: timestamppb.New(groupMember.UpdatedAt),
		},
	}, nil

}
func (g *groupGRPCServer) DeleteGroupMember(ctx context.Context, request *v1.DeleteGroupMemberRequest) (*v1.DeleteGroupMemberResponse, error) {
	slog.Info("DELETE: DeleteGroupMember")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.ucGroupMember.DeleteGroupMember(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
	}
	return &v1.DeleteGroupMemberResponse{
		Deleted: deleted,
	}, nil
}
func (g *groupGRPCServer) UpdateGroupMember(ctx context.Context, request *v1.UpdateGroupMemberRequest) (*v1.UpdateGroupMemberResponse, error) {
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
	res := &v1.UpdateGroupMemberResponse{
		GroupMember: &v1.GroupMember{
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
func (g *groupGRPCServer) GetGroups(ctx context.Context, request *v1.GetGroupsRequest) (*v1.GetGroupsResponse, error) {
	slog.Info("GET: GetGroupsByUserId")

	groups, err := g.ucGroup.GetGroups(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetGroupsByUserId failed")
	}
	return &v1.GetGroupsResponse{
		Groups: lo.Map(groups, func(group *domain.Group, _ int) *v1.Group {
			return entityToProtobuf(group)
		}),
	}, nil
}
func (g *groupGRPCServer) GetGroupsByUserId(ctx context.Context, request *v1.GetGroupsByUserIdRequest) (*v1.GetGroupsByUserIdResponse, error) {
	slog.Info("GET: GetGroupsByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groups, err := g.ucGroup.GetGroupsByUserId(ctx, userId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetGroupsByUserId failed")
	}
	return &v1.GetGroupsByUserIdResponse{
		Groups: lo.Map(groups, func(group *domain.Group, _ int) *v1.Group {
			return entityToProtobuf(group)
		}),
	}, nil
}

func (g *groupGRPCServer) GetGroupIdsByUserId(ctx context.Context, request *v1.GetGroupIdsByUserIdRequest) (*v1.GetGroupIdsByUserIdResponse, error) {
	slog.Info("GET: GetAllGroupByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groupIds, err := g.ucGroup.GetGroupIdsByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetGroupIdsByUserId failed")
	}
	return &v1.GetGroupIdsByUserIdResponse{
		GroupIds: groupIds,
	}, nil
}

func (g *groupGRPCServer) CreateGroup(ctx context.Context, request *v1.CreateGroupRequest) (*v1.CreateGroupResponse, error) {
	slog.Info("POST: CreateGroup")
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	model := domain.Group{
		Name:        request.Group.Name,
		Description: request.Group.Description,
		Status:      request.Group.Status,
		UserID:      payloadUser.ID,
	}
	group, err := g.ucGroup.CreateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.CreateGroup failed")
	}
	if request.Group.UserIds != nil {
		for _, userId := range request.Group.UserIds {
			userIdParsed, err := uuid.Parse(userId)
			if err != nil {
				return nil, errors.Wrap(err, "ucGroup.CreateGroup userIds parse failed")
			}
			g.ucGroupMember.CreateGroupMember(ctx, &domain.GroupMember{
				ID:      uuid.New(),
				GroupID: group.ID,
				UserID:  userIdParsed,
				Role:    int32(sharedkernel.USER),
			})
		}
	}
	res := &v1.CreateGroupResponse{
		Group: entityToProtobuf(group),
	}
	return res, nil
}

func (g *groupGRPCServer) GetGroup(ctx context.Context, request *v1.GetGroupRequest) (*v1.GetGroupResponse, error) {
	slog.Info("GET: GetGroup")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	group, err := g.ucGroup.GetGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
	}
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	countMembers, err := g.ucGroupMember.CountGroupMembers(ctx, group.ID)
	if err != nil {
		countMembers = 0
	}
	roleMember, err := g.ucGroupMember.GetRoleOfGroupMember(ctx, group.ID, payloadUser.ID)
	if err != nil {
		roleMember = 0
	}
	// isOwner , err := g
	return &v1.GetGroupResponse{
		Group:             entityToProtobuf(group),
		CountGroupMembers: countMembers,
		RoleGroupMember:   sharedkernel.RoleGroupMember(roleMember).String(),
	}, nil
}

func (g *groupGRPCServer) DeleteGroup(ctx context.Context, request *v1.DeleteGroupRequest) (*v1.DeleteGroupResponse, error) {
	slog.Info("DELETE: DeleteGroup")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.ucGroup.DeleteGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.DeleteGroup failed")
	}
	return &v1.DeleteGroupResponse{
		Deleted: deleted,
	}, nil
}

func (g *groupGRPCServer) UpdateGroup(ctx context.Context, request *v1.UpdateGroupRequest) (*v1.UpdateGroupResponse, error) {
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
		ProfileUrl:  request.Group.ProfileUrl,
	}

	group, err := g.ucGroup.UpdateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.UpdateGroup failed")
	}
	res := &v1.UpdateGroupResponse{
		Group: entityToProtobuf(group),
	}
	return res, nil
}

func entityToProtobuf(group *domain.Group) *v1.Group {
	return &v1.Group{
		Id:          group.ID.String(),
		Name:        group.Name,
		Description: group.Description,
		UserId:      group.UserID.String(),
		Status:      group.Status,
		ProfileUrl:  group.ProfileUrl,
		CreatedAt:   timestamppb.New(group.CreatedAt),
		UpdatedAt:   timestamppb.New(group.UpdatedAt),
	}
}
