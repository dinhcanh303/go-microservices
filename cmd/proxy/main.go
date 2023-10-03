package main

import (
	"context"
	"go-microservices/cmd/proxy/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func newGateway(ctx context.Context, cfg *config.Config, opts []runtime)
