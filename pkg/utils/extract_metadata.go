package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
		if len(userValues) == 5 {
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
func ExtractHeaderUser(ctx echo.Context) (*token.Payload, error) {
	headerUser := ctx.Request().Header.Get("Grpc-Metadata-X-Auth-User")
	payload := &token.Payload{}
	if headerUser != "" {
		userValues := strings.Split(headerUser, ",")
		if len(userValues) == 5 {
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
		return nil, errors.New("header not found header forward")
	}
	return payload, nil
}
func ExtractMetadataKeyStore(ctx context.Context) (*domain.Key, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	slog.Info("CONTEXT::", md)
	keyStore := &domain.Key{}
	values := md.Get(constant.KeyStore)
	values2 := md.Get(constant.KeyStoreUsed)
	if len(values) > 0 && len(values2) > 0 {
		keyStoreValues := strings.Split(values[0], ",")
		keyStoreUsedValues := strings.Split(values2[0], ",")
		if len(keyStoreValues) == 5 {
			id, _ := strconv.ParseInt(keyStoreValues[0], 10, 64)
			keyStore.ID = id
			userId, err := uuid.Parse(keyStoreValues[1])
			if err != nil {
				return nil, err
			}
			keyStore.UserID = userId
			keyStore.PublicKey = keyStoreValues[2]
			keyStore.PrivateKey = keyStoreValues[3]
			keyStore.RefreshToken = keyStoreValues[4]
		}
		if len(keyStoreUsedValues) != 0 {
			jsonBytes, _ := json.Marshal(keyStoreUsedValues)
			keyStore.RefreshTokensUsed = json.RawMessage(jsonBytes)
		}
	} else {
		return nil, errors.New("context not found header forward")
	}
	return keyStore, nil
}
