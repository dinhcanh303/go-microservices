package utils

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func ExtractMetadataUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	slog.Info("CONTEXT::", md)
	payload := &token.Payload{}
	if values := md.Get(constant.User); len(values) > 0 {
		userValues := strings.Split(values[0], ",")
		if len(values) == 5 {
			userId, err := uuid.Parse(userValues[0])
			if err != nil {
				return nil, err
			}
			payload.ID = userId
			payload.Email = userValues[1]
			payload.FullName = userValues[2]
			payload.Role = userValues[3]
			payload.AvatarUrl = userValues[4]
		}
	} else {
		return nil, errors.New("context not found header forward")
	}
	return payload, nil
}
