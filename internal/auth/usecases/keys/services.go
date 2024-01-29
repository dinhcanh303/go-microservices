package keys

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo  KeyRepo
	redis redis.RedisEngine
}

var _ UseCase = (*service)(nil)

func NewUseCase(repo KeyRepo,
	redis redis.RedisEngine) UseCase {
	return &service{
		repo:  repo,
		redis: redis,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)

// CreateKeyToken implements UseCase.
func (s *service) CreateKeyToken(ctx context.Context, key *domain.Key) (*domain.Key, error) {
	userID := key.UserID
	foundKeyToken, err := s.repo.FindKeyByUserID(ctx, userID)
	if err != nil {
		slog.Warn("CreateKeyToken::", err)
		keyToken, err := s.repo.CreateKey(ctx, key)
		if err != nil {
			return nil, errors.Wrap(err, "Create Key Token failed")
		}
		//Del cache
		err = s.redis.Invalidate(constant.CacheKeyTokenUser + userID.String())
		if err != nil {
			slog.Error("Invalidate cache key failed")
		}
		return keyToken, nil
	}

	keyToken, err := s.repo.UpdateKeyByUserID(ctx, &domain.Key{
		UserID:            foundKeyToken.UserID,
		PublicKey:         key.PublicKey,
		PrivateKey:        key.PrivateKey,
		RefreshToken:      key.RefreshToken,
		RefreshTokensUsed: key.RefreshTokensUsed,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Update Key Token failed")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheKeyTokenUser + userID.String())
	if err != nil {
		slog.Error("Invalidate cache key failed")
	}
	return keyToken, nil
}

// DeleteKeyByID implements UseCase.
func (s *service) DeleteKeyByID(ctx context.Context, id int64) error {
	err := s.repo.DeleteKeyByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "service.DeleteKeyByID failed")
	}

	return nil
}

// DeleteKeyByUserID implements UseCase.
func (s *service) DeleteKeyByUserID(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteKeyByUserID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "service.DeleteKeyByUserID failed")
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheKeyTokenUser + userID.String())
	if err != nil {
		slog.Error("Invalidate cache key failed")
	}
	return nil
}

// FindKeyByRefreshToken implements UseCase.
func (s *service) FindKeyByRefreshToken(ctx context.Context, refreshToken string) (*domain.Key, error) {
	key, err := s.repo.FindKeyByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "service.FindKeyByRefreshTokenUsed failed")
	}
	return key, nil
}

// FindKeyByRefreshTokenUsed implements UseCase.
func (s *service) FindKeyByRefreshTokenUsed(ctx context.Context, refreshToken string) (*domain.Key, error) {
	key, err := s.repo.FindKeyByRefreshTokenUsed(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "service.FindKeyByRefreshTokenUsed failed")
	}
	return key, nil
}

// FindKeyByUserID implements UseCase.
func (s *service) FindKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.Key, error) {
	var key *domain.Key
	keyCache := constant.CacheKeyTokenUser + userID.String()
	err := utils.HandleHitCache(key, s.redis, keyCache)
	if err != nil {
		key, err = s.repo.FindKeyByUserID(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, "service.FindKeyByUserID failed")
		}
		err = s.redis.Set(keyCache, key)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return key, nil
}
