package keys

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type service struct {
	repo KeyRepo
}

// CreateKeyToken implements UseCase.
func (s *service) CreateKeyToken(ctx context.Context, key *domain.Key) (*domain.Key, error) {
	userID := key.UserID
	foundKeyToken, err := s.repo.FindKeyByUserID(ctx, userID)
	slog.Warn("Service::", err)
	if err != nil {
		keyToken, err := s.repo.CreateKey(ctx, key)
		if err != nil {
			return nil, errors.Wrap(err, "Create Key Token failed")
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
	slog.Warn("Service::", err)
	if err != nil {
		return nil, errors.Wrap(err, "Update Key Token failed")
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
	key, err := s.repo.FindKeyByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "service.FindKeyByUserID failed")
	}
	return key, nil
}

var _ UseCase = (*service)(nil)

func NewUseCase(repo KeyRepo) UseCase {
	return &service{
		repo: repo,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)
