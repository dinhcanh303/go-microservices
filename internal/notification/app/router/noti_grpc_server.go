package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type notiGRPCServer struct {
	gen.UnimplementedNotiServiceServer
	uc notifications.UseCase
}

var _ gen.NotiServiceServer = (*notiGRPCServer)(nil)

var NotiGRPCServerSet = wire.NewSet(NewNotiGRPCServer)

func NewNotiGRPCServer(
	grpcServer *grpc.Server,
	uc notifications.UseCase,
) gen.NotiServiceServer {
	svc := notiGRPCServer{
		uc: uc,
	}
	gen.RegisterNotiServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (n *notiGRPCServer) GetNotifications(ctx context.Context, request *gen.GetNotificationsRequest) (*gen.GetNotificationsResponse, error) {

	panic("not")
}
