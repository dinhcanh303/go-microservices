package router

import (
	"github.com/dinhcanh303/go-microservices/cmd/group/config"
	"github.com/dinhcanh303/go-microservices/internal/group/usecases/groups"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type groupGRPCServer struct {
	gen
	cfg *config.Config
	uc  groups.UseCase
}

func NewGRPCCounterServer() {

}
