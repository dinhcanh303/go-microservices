package handlers

import (
	"context"

	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/upload/app/request"
	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/pkg/echo/responses"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uc uploads.UseCase
}

func NewUploadHandler(uc uploads.UseCase) *UploadHandler {
	return &UploadHandler{
		uc: uc,
	}
}

var UploadServerSet = wire.NewSet(NewUploadHandler)

// DeleteAttachment implements uploads.UseCase.
func (s *UploadHandler) DeleteAttachment(ctx echo.Context) error {
	deleteAttachment := new(request.DeleteAttachmentRequest)
	if err := ctx.Bind(deleteAttachment); err != nil {
		return err
	}
	attachment, err := s.uc.DeleteAttachment(context.Background(), deleteAttachment.AttachmentID)
	if err != nil {
		return err
	}
	return responses.Response(ctx, http.StatusOK, attachment)
}

// GetAttachment implements uploads.UseCase.
func (s *UploadHandler) GetAttachment(ctx echo.Context) error {
	getAttachment := new(request.GetAttachmentRequest)
	if err := ctx.Bind(getAttachment); err != nil {
		return err
	}
	attachment, err := s.uc.GetAttachment(context.Background(), getAttachment.AttachmentID)
	if err != nil {
		return err
	}
	return responses.Response(ctx, http.StatusOK, attachment)

}

// UpdateAttachment implements uploads.UseCase.
func (s *UploadHandler) UpdateAttachment(ctx echo.Context) error {
	updateAttachment := new(request.UpdateAttachmentRequest)
	if err := ctx.Bind(updateAttachment); err != nil {
		return err
	}
	attachment, err := s.uc.UpdateAttachment(context.Background(), &domain.Attachment{
		AttachableType: updateAttachment.AttachableType,
		AttachableID:   updateAttachment.AttachableID,
	})
	if err != nil {
		return err
	}
	return responses.Response(ctx, http.StatusOK, attachment)
}

// UploadFile implements uploads.UseCase.
func (s *UploadHandler) UploadFile(ctx echo.Context) error {
	attachments, err := s.uc.UploadFile(ctx)
	if err != nil {
		return err
	}
	return responses.Response(ctx, http.StatusOK, attachments)
}
