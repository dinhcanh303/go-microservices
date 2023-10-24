package router

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type uploadGRPCServer struct {
	gen.UnimplementedUploadServiceServer
	cfg *config.Config
	uc  uploads.UseCase
}

var _ gen.UploadServiceServer = (*uploadGRPCServer)(nil)
