package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sqlc-dev/pqtype"
)

type keyRepo struct {
	pg postgres.DBEngine
}

// CreateKey implements keys.KeyRepo.
func (rp *keyRepo) CreateKey(ctx context.Context, key *domain.Key) (*domain.Key, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.CreateKey db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.CreateKey(ctx, postgresql.CreateKeyParams{
		UserID:     key.UserID,
		PublicKey:  key.PublicKey,
		PrivateKey: key.PrivateKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.CreateKey failed")
	}
	return &domain.Key{
		ID:                result.ID,
		PublicKey:         result.PublicKey,
		PrivateKey:        result.PrivateKey,
		UserID:            result.UserID,
		RefreshToken:      result.RefreshToken.String,
		RefreshTokensUsed: result.RefreshTokensUsed.RawMessage,
		CreatedAt:         result.CreatedAt,
		UpdatedAt:         result.UpdatedAt,
	}, tx.Commit()
}

// UpdateKeyByUserID implements keys.KeyRepo.
func (rp *keyRepo) UpdateKeyByUserID(ctx context.Context, key *domain.Key) (*domain.Key, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.UpdateKey db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.UpdateKeyByUserID(ctx, postgresql.UpdateKeyByUserIDParams{
		UserID: key.UserID,
		PublicKey: sql.NullString{
			String: key.PublicKey,
			Valid:  key.PublicKey != "",
		},
		PrivateKey: sql.NullString{
			String: key.PrivateKey,
			Valid:  key.PrivateKey != "",
		},
		RefreshToken: sql.NullString{
			String: key.RefreshToken,
			Valid:  key.RefreshToken != "",
		},
		RefreshTokensUsed: pqtype.NullRawMessage{
			RawMessage: key.RefreshTokensUsed,
			Valid:      key.RefreshTokensUsed != nil,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.UpdateKey failed")
	}
	return &domain.Key{
		ID:                result.ID,
		PublicKey:         result.PublicKey,
		PrivateKey:        result.PrivateKey,
		UserID:            result.UserID,
		RefreshToken:      result.RefreshToken.String,
		RefreshTokensUsed: result.RefreshTokensUsed.RawMessage,
		CreatedAt:         result.CreatedAt,
		UpdatedAt:         result.UpdatedAt,
	}, tx.Commit()
}

// DeleteKeyByID implements keys.KeyRepo.
func (rp *keyRepo) DeleteKeyByID(ctx context.Context, id int64) error {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "keyRepo.UpdateKey db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteKeyByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "keyRepo.UpdateKey failed")
	}
	return tx.Commit()
}

// DeleteKeyByUserID implements keys.KeyRepo.
func (rp *keyRepo) DeleteKeyByUserID(ctx context.Context, userID uuid.UUID) error {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "keyRepo.UpdateKey db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteKeyByUserID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "keyRepo.UpdateKey failed")
	}
	return tx.Commit()
}

// FindKeyByRefreshToken implements keys.KeyRepo.
func (rp *keyRepo) FindKeyByRefreshToken(ctx context.Context, refreshToken string) (*domain.Key, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)

	key, err := querier.FindKeyByRefreshToken(ctx, sql.NullString{
		String: refreshToken,
		Valid:  refreshToken != "",
	})
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.UpdateKey failed")
	}
	return &domain.Key{
		ID:                key.ID,
		UserID:            key.UserID,
		PublicKey:         key.PublicKey,
		PrivateKey:        key.PrivateKey,
		RefreshToken:      key.RefreshToken.String,
		RefreshTokensUsed: key.RefreshTokensUsed.RawMessage,
		CreatedAt:         key.CreatedAt,
		UpdatedAt:         key.UpdatedAt,
	}, nil
}

// FindKeyByRefreshTokenUsed implements keys.KeyRepo.
func (*keyRepo) FindKeyByRefreshTokenUsed(ctx context.Context, refreshToken string) (*domain.Key, error) {
	panic("unimplemented")
}

// FindKeyByUserID implements keys.KeyRepo.
func (rp *keyRepo) FindKeyByUserID(ctx context.Context, userID uuid.UUID) (*domain.Key, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)

	key, err := querier.FindKeyByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "keyRepo.UpdateKey failed")
	}
	return &domain.Key{
		ID:                key.ID,
		UserID:            key.UserID,
		PublicKey:         key.PublicKey,
		PrivateKey:        key.PrivateKey,
		RefreshToken:      key.RefreshToken.String,
		RefreshTokensUsed: key.RefreshTokensUsed.RawMessage,
		CreatedAt:         key.CreatedAt,
		UpdatedAt:         key.UpdatedAt,
	}, nil
}

var _ keys.KeyRepo = (*keyRepo)(nil)

func NewKeyRepo(pg postgres.DBEngine) keys.KeyRepo {
	return &keyRepo{pg: pg}
}

var KeyRepoSet = wire.NewSet(NewKeyRepo)
