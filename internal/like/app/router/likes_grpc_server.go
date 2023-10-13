package router

import (
	"github.com/dinhcanh303/go-microservices/cmd/like/config"
	"github.com/dinhcanh303/go-microservices/internal/like/usecases/likes"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
