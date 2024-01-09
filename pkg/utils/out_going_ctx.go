package utils

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

func OutgoingContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	ctxBackground := metadata.NewOutgoingContext(ctx, md)
	return ctxBackground, nil
}
