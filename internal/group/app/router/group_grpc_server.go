package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	gen "github.com/dinhcanh303/go-microservices/gen/go"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type groupGRPCServer struct {
	gen.UnimplementedGroupServiceServer
	cfg *config.Config
	uc  groups.UseCase
}

var _ gen.UnimplementedGroupServiceServer = (*groupGRPCServer)(nil)

var GroupGRPCServerSet = wire.NewSet(NewGRPCCounterServer)

func NewGRPCCounterServer(
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
	id, err := uuid.Parse(request.Group.Uuid)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse failed")
	}

}
