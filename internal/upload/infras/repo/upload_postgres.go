package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/internal/upload/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type attachmentRepo struct {
	pg postgres.DBEngine
}

// Create implements uploads.AttachmentRepo.
func (rp *attachmentRepo) Create(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Create(ctx, postgresql.CreateParams{
		ID:        attachment.ID,
		UserID:    attachment.UserID,
		Filename:  attachment.FileName,
		Extension: attachment.Extension,
		MimeType: sql.NullString{
			String: attachment.MimeType,
			Valid:  attachment.MimeType != "",
		},
		VersionID: sql.NullString{
			String: attachment.VersionID,
			Valid:  attachment.VersionID != "",
		},
		Url: attachment.URL,
		UrlThumbnail: sql.NullString{
			String: attachment.URLThumbnail,
			Valid:  attachment.URLThumbnail != "",
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return &domain.Attachment{
		ID:             result.ID,
		AttachableType: result.AttachableType.String,
		AttachableID:   result.AttachableID.UUID,
		UserID:         result.UserID,
		FileName:       result.Filename,
		Extension:      result.Extension,
		MimeType:       result.MimeType.String,
		Folder:         result.Folder.String,
		VersionID:      result.VersionID.String,
		URL:            result.Url,
		URLThumbnail:   result.UrlThumbnail.String,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
	}, tx.Commit()
}

// Delete implements uploads.AttachmentRepo.
func (rp *attachmentRepo) Delete(ctx context.Context, attachmentId uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "postRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.Delete(ctx, attachmentId)
	if err != nil {
		return false, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return true, tx.Commit()
}

// Get implements uploads.AttachmentRepo.
func (rp *attachmentRepo) Get(ctx context.Context, attachmentId uuid.UUID) (*domain.Attachment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	result, err := querier.Get(ctx, attachmentId)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, attachmentId) failed")
	}
	return &domain.Attachment{
		ID:             result.ID,
		AttachableType: result.AttachableType.String,
		AttachableID:   result.AttachableID.UUID,
		UserID:         result.UserID,
		FileName:       result.Filename,
		Extension:      result.Extension,
		MimeType:       result.MimeType.String,
		Folder:         result.Folder.String,
		VersionID:      result.VersionID.String,
		URL:            result.Url,
		URLThumbnail:   result.UrlThumbnail.String,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
	}, nil
}

// Update implements uploads.AttachmentRepo.
func (rp *attachmentRepo) Update(ctx context.Context, attachment *domain.Attachment) (*domain.Attachment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create db failed")
	}
	qtx := querier.WithTx(tx)
	result, err := qtx.Update(ctx, postgresql.UpdateParams{
		ID: attachment.ID,
		AttachableType: sql.NullString{
			String: attachment.AttachableType,
			Valid:  attachment.AttachableType != "",
		},
		AttachableID: uuid.NullUUID{
			UUID:  attachment.AttachableID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Create(ctx, postgresql.CreateParams) failed")
	}

	return &domain.Attachment{
		ID:             result.ID,
		AttachableType: result.AttachableType.String,
		AttachableID:   result.AttachableID.UUID,
		UserID:         result.UserID,
		FileName:       result.Filename,
		Extension:      result.Extension,
		MimeType:       result.MimeType.String,
		Folder:         result.Folder.String,
		VersionID:      result.VersionID.String,
		URL:            result.Url,
		URLThumbnail:   result.UrlThumbnail.String,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
	}, tx.Commit()
}

func NewAttachmentRepo(pg postgres.DBEngine) uploads.AttachmentRepo {
	return &attachmentRepo{pg: pg}
}

var _ uploads.AttachmentRepo = (*attachmentRepo)(nil)

var RepositoryUploadSet = wire.NewSet(NewAttachmentRepo)
