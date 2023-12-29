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
	"github.com/samber/lo"
)

type attachmentRepo struct {
	pg postgres.DBEngine
}

// GetLastAttachmentByType implements uploads.AttachmentRepo.
func (rp *attachmentRepo) GetLastAttachmentByType(ctx context.Context, attachableType string, attachableId uuid.UUID) (*domain.Attachment, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	result, err := querier.GetLastAttachmentByType(ctx, postgresql.GetLastAttachmentByTypeParams{
		AttachableType: sql.NullString{
			String: attachableType,
			Valid:  attachableType != "",
		},
		AttachableID: uuid.NullUUID{
			UUID:  attachableId,
			Valid: true,
		},
	})
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
		URL:            result.Url,
		URLThumbnail:   result.UrlThumbnail.String,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
	}, nil
}

// GetAttachmentsByType implements uploads.AttachmentRepo.
func (rp *attachmentRepo) GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domain.Attachment, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	results, err := querier.GetAttachmentsByType(ctx, postgresql.GetAttachmentsByTypeParams{
		AttachableType: sql.NullString{
			String: attachableType,
			Valid:  attachableType != "",
		},
		AttachableID: uuid.NullUUID{
			UUID:  attachableId,
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, attachmentId) failed")
	}
	return lo.Map(results, func(item postgresql.UploadAttachment, _ int) *domain.Attachment {
		return &domain.Attachment{
			ID:             item.ID,
			AttachableType: item.AttachableType.String,
			AttachableID:   item.AttachableID.UUID,
			UserID:         item.UserID,
			FileName:       item.Filename,
			Extension:      item.Extension,
			MimeType:       item.MimeType.String,
			Folder:         item.Folder.String,
			URL:            item.Url,
			URLThumbnail:   item.UrlThumbnail.String,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}), nil
}

// UpdateByIds implements uploads.AttachmentRepo.
func (rp *attachmentRepo) UpdateByIds(ctx context.Context, attachmentIds []uuid.UUID, attachment *domain.Attachment) ([]*domain.Attachment, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "attachmentRepo.UpdateByIds db failed")
	}
	qtx := querier.WithTx(tx)
	results, err := qtx.UpdateByIds(ctx, postgresql.UpdateByIdsParams{
		AttachableType: sql.NullString{
			String: attachment.AttachableType,
			Valid:  attachment.AttachableType != "",
		},
		AttachableID: uuid.NullUUID{
			UUID:  attachment.AttachableID,
			Valid: true,
		},
		Column1: attachmentIds,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.UpdateByIds(ctx, postgresql.UpdateByIdsParams) failed")
	}
	return lo.Map(results, func(item postgresql.UploadAttachment, _ int) *domain.Attachment {
		return &domain.Attachment{
			ID:             item.ID,
			AttachableType: item.AttachableType.String,
			AttachableID:   item.AttachableID.UUID,
			UserID:         item.UserID,
			FileName:       item.Filename,
			Extension:      item.Extension,
			MimeType:       item.MimeType.String,
			Folder:         item.Folder.String,
			URL:            item.Url,
			URLThumbnail:   item.UrlThumbnail.String,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}), tx.Commit()
}

// DeleteByIds implements uploads.AttachmentRepo.
func (rp *attachmentRepo) DeleteByIds(ctx context.Context, attachmentIds []uuid.UUID) (bool, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)
	tx, err := db.Begin()
	if err != nil {
		return false, errors.Wrap(err, "attachmentRepo.DeleteByIds db failed")
	}
	qtx := querier.WithTx(tx)
	err = qtx.DeleteByIds(ctx, attachmentIds)
	if err != nil {
		return false, errors.Wrap(err, "attachmentRepo.DeleteByIds db failed")
	}
	return true, tx.Commit()
}

// GetByIds implements uploads.AttachmentRepo.
func (rp *attachmentRepo) GetByIds(ctx context.Context, attachmentIds []uuid.UUID) ([]*domain.Attachment, error) {
	db := rp.pg.GetDBRead()
	querier := postgresql.New(db)
	results, err := querier.GetByIds(ctx, attachmentIds)
	if err != nil {
		return nil, errors.Wrap(err, "qtx.Get(ctx, attachmentId) failed")
	}
	return lo.Map(results, func(item postgresql.UploadAttachment, _ int) *domain.Attachment {
		return &domain.Attachment{
			ID:             item.ID,
			AttachableType: item.AttachableType.String,
			AttachableID:   item.AttachableID.UUID,
			UserID:         item.UserID,
			FileName:       item.Filename,
			Extension:      item.Extension,
			MimeType:       item.MimeType.String,
			Folder:         item.Folder.String,
			URL:            item.Url,
			URLThumbnail:   item.UrlThumbnail.String,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}), nil
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
		Folder: sql.NullString{
			String: attachment.Folder,
			Valid:  attachment.Folder != "",
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
	db := rp.pg.GetDBRead()
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
		return nil, errors.Wrap(err, "attachmentRepo.Create db failed")
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
