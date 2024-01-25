package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/domain"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groupmembers"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
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

type groupGRPCServer struct {
	gen.UnimplementedGroupServiceServer
	cfg               *config.Config
	ucGroup           groups.UseCase
	ucGroupMember     groupmembers.UseCase
	authDomainService domain.AuthDomainService
	redis             redis.RedisEngine
}

var _ gen.GroupServiceServer = (*groupGRPCServer)(nil)

var GroupGRPCServerSet = wire.NewSet(NewGRPCGroupServer)

func NewGRPCGroupServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	ucGroup groups.UseCase,
	ucGroupMember groupmembers.UseCase,
	redis redis.RedisEngine,
	authDomainService domain.AuthDomainService,
) gen.GroupServiceServer {
	svc := groupGRPCServer{
		cfg:               cfg,
		ucGroup:           ucGroup,
		ucGroupMember:     ucGroupMember,
		redis:             redis,
		authDomainService: authDomainService,
	}
	gen.RegisterGroupServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *groupGRPCServer) GetGroupMembers(ctx context.Context, request *gen.GetGroupMembersRequest) (*gen.GetGroupMembersResponse, error) {
	slog.Info("GET: GetAllGroupMembers")
	groupId, err := uuid.Parse(request.GroupId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}
	groupMembers, err := g.ucGroupMember.GetGroupMembers(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroupMember.GetAllGroupMembers failed")
	}
	return &gen.GetGroupMembersResponse{
		GroupMembers: lo.Map(groupMembers, func(groupMember *domain.GroupMember, _ int) *gen.GroupMemberMetadata {
			user, err := g.authDomainService.GetProfile(ctx, groupMember.UserID)
			if err != nil {
				user = &gen.GetProfileResponse{}
			}
			return &gen.GroupMemberMetadata{
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
func (g *groupGRPCServer) CreateGroupMember(ctx context.Context, request *gen.CreateGroupMemberRequest) (*gen.CreateGroupMemberResponse, error) {
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
	return &gen.CreateGroupMemberResponse{
		GroupMember: &gen.GroupMember{
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
		GroupMember: &gen.GroupMember{
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

func (g *groupGRPCServer) GetGroupsByUserId(ctx context.Context, request *gen.GetGroupsByUserIdRequest) (*gen.GetGroupsByUserIdResponse, error) {
	slog.Info("GET: GetGroupsByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groups, err := g.ucGroup.GetGroupsByUserId(ctx, userId, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetGroupsByUserId failed")
	}
	return &gen.GetGroupsByUserIdResponse{
		Groups: lo.Map(groups, func(group *domain.Group, _ int) *gen.Group {
			return &gen.Group{
				Id:          group.ID.String(),
				Name:        group.Name,
				Description: group.Description,
				Status:      group.Status,
				UserId:      group.UserID.String(),
				ProfileUrl:  group.ProfileUrl,
				CreatedAt:   timestamppb.New(group.CreatedAt),
				UpdatedAt:   timestamppb.New(group.UpdatedAt),
			}
		}),
	}, nil
}

func (g *groupGRPCServer) GetGroupIdsByUserId(ctx context.Context, request *gen.GetGroupIdsByUserIdRequest) (*gen.GetGroupIdsByUserIdResponse, error) {
	slog.Info("GET: GetAllGroupByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse(request.UserId) failed")
	}
	groupIds, err := g.ucGroup.GetGroupIdsByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "userId.GetGroupIdsByUserId failed")
	}
	return &gen.GetGroupIdsByUserIdResponse{
		GroupIds: groupIds,
	}, nil
}

func (g *groupGRPCServer) CreateGroup(ctx context.Context, request *gen.CreateGroupRequest) (*gen.CreateGroupResponse, error) {
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
	//Del cache
	err = g.redis.Invalidate(group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
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
				Role:    constant.MEMBER,
			})
		}
	}
	res := &gen.CreateGroupResponse{
		Group: &gen.Group{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			ProfileUrl:  group.ProfileUrl,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
	}
	return res, nil
}

func (g *groupGRPCServer) GetGroup(ctx context.Context, request *gen.GetGroupRequest) (*gen.GetGroupResponse, error) {
	slog.Info("GET: GetGroup")
	group := &domain.Group{}
	groupId := request.Id
	err := utils.HandleHitCache(group, g.redis, groupId)
	if err != nil {
		slog.Info("MISS_CACHE", err)
		id, err := uuid.Parse(groupId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse id")
		}
		group, err = g.ucGroup.GetGroup(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "ucGroup.GetGroup failed")
		}
		err = g.redis.Set(group.ID.String(), group, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	countMembers, err := g.ucGroupMember.CountGroupMembers(ctx, group.ID)
	if err != nil {
		countMembers = 0
	}
	isMember, err := g.ucGroupMember.CheckGroupMember(ctx, group.ID, payloadUser.ID)
	if err != nil {
		isMember = false
	}
	return &gen.GetGroupResponse{
		Group: &gen.Group{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			ProfileUrl:  group.ProfileUrl,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
		CountGroupMembers: countMembers,
		IsMember:          isMember,
	}, nil
}

func (g *groupGRPCServer) DeleteGroup(ctx context.Context, request *gen.DeleteGroupRequest) (*gen.DeleteGroupResponse, error) {
	slog.Info("DELETE: DeleteGroup")
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse id")
	}
	deleted, err := g.ucGroup.DeleteGroup(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.DeleteGroup failed")
	}

	//Del cache
	err = g.redis.Invalidate(request.Id)
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", request.Id)
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
		ProfileUrl:  request.Group.ProfileUrl,
	}

	group, err := g.ucGroup.UpdateGroup(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "ucGroup.UpdateGroup failed")
	}

	//Del cache
	err = g.redis.Invalidate(group.ID.String())
	if err != nil {
		slog.Error("Invalidate cache group failed : ID", group.ID.String())
	}
	res := &gen.UpdateGroupResponse{
		Group: &gen.Group{
			Id:          group.ID.String(),
			Name:        group.Name,
			Description: group.Description,
			UserId:      group.UserID.String(),
			Status:      group.Status,
			ProfileUrl:  group.ProfileUrl,
			CreatedAt:   timestamppb.New(group.CreatedAt),
			UpdatedAt:   timestamppb.New(group.UpdatedAt),
		},
	}
	return res, nil
}
