package router

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/attachment"
)

type uploadGRPCServer struct {
	cfg *config.Config
	uc  attachment.UseCase
}

var _ gen.UploadServiceServer = (*uploadGRPCServer)(nil)
