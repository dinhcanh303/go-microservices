package handlers

import (
	"context"
	"log/slog"
	"slices"

	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/upload/app/request"
	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/echo/responses"
	"github.com/google/uuid"
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

var UploadHandlerSet = wire.NewSet(NewUploadHandler)

// DeleteAttachment implements uploads.UseCase.
func (s *UploadHandler) DeleteAttachmentsByIds(ctx echo.Context) error {
	deleteAttachments := new(request.DeleteAttachmentsByIdsRequest)
	if err := ctx.Bind(deleteAttachments); err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Bind failed", err))
	}
	attachmentIds := make([]uuid.UUID, 0)
	for _, id := range deleteAttachments.AttachmentIds {
		attachmentId, err := uuid.Parse(id)
		if err != nil {
			return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
		}
		attachmentIds = append(attachmentIds, attachmentId)
	}
	deleted, err := s.uc.DeleteAttachmentsByIds(context.Background(), attachmentIds)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler DeleteAttachmentsByIds failed", err))
	}
	return responses.Response(ctx, http.StatusOK, deleted)
}

// DeleteAttachment implements uploads.UseCase.
func (s *UploadHandler) DeleteAttachment(ctx echo.Context) error {
	id := ctx.Param("id")
	attachmentId, err := uuid.Parse(id)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
	}
	deleted, err := s.uc.DeleteAttachment(context.Background(), attachmentId)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler DeleteAttachment failed", err))
	}
	return responses.Response(ctx, http.StatusOK, deleted)
}

// GetAttachment implements uploads.UseCase.
func (s *UploadHandler) GetAttachment(ctx echo.Context) error {
	id := ctx.Param("id")
	attachmentId, err := uuid.Parse(id)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
	}
	attachment, err := s.uc.GetAttachment(context.Background(), attachmentId)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler UploadFile failed", err))
	}
	return responses.Response(ctx, http.StatusOK, attachment)

}

// UpdateAttachment implements uploads.UseCase.
func (s *UploadHandler) UpdateAttachment(ctx echo.Context) error {
	slog.Info("PUT: UpdateAttachment")
	id := ctx.Param("id")
	updateAttachment := new(request.UpdateAttachmentRequest)
	if err := ctx.Bind(updateAttachment); err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Bind failed", err))
	}
	attachableType := updateAttachment.AttachableType
	if !checkAttachmentTypeByCondition(attachableType, nil) {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Attachment Type isn't correct", nil))
	}
	attachmentId, err := uuid.Parse(id)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
	}
	attachableId, err := uuid.Parse(updateAttachment.AttachableID)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
	}
	entityUpload := updateAttachment.EntityUpload
	attachment, err := s.uc.UpdateAttachment(context.Background(), &domain.Attachment{
		ID:             attachmentId,
		AttachableType: attachableType,
		AttachableID:   attachableId,
		EntityUpload:   entityUpload,
	})
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler UpdateAttachment failed", err))
	}
	return responses.Response(ctx, http.StatusOK, attachment)
}

// UpdateAttachmentsByIds implements uploads.UseCase.
func (s *UploadHandler) UpdateAttachmentsByIds(ctx echo.Context) error {
	slog.Info("POST: UpdateAttachmentsByIds")
	updateAttachment := new(request.UpdateAttachmentsByIdsRequest)
	if err := ctx.Bind(updateAttachment); err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Bind failed", err))
	}
	attachableType := updateAttachment.AttachableType
	if !checkAttachmentTypeByCondition(attachableType, nil) {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Attachment Type isn't correct", nil))
	}
	attachmentIds := make([]uuid.UUID, 0)
	for _, id := range updateAttachment.AttachmentIds {
		attachmentId, err := uuid.Parse(id)
		if err != nil {
			return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
		}
		attachmentIds = append(attachmentIds, attachmentId)
	}

	attachableId, err := uuid.Parse(updateAttachment.AttachableID)
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Parse failed", err))
	}
	entityUpload := updateAttachment.EntityUpload
	attachment, err := s.uc.UpdateAttachmentsByIds(context.Background(), attachmentIds, &domain.Attachment{
		AttachableType: attachableType,
		AttachableID:   attachableId,
		EntityUpload:   entityUpload,
	})
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler UpdateAttachmentsByIds failed", err))
	}
	return responses.Response(ctx, http.StatusOK, attachment)
}

// UploadFile implements uploads.UseCase.
func (s *UploadHandler) UploadFile(ctx echo.Context) error {
	slog.Info("POST: UploadFile")
	attachments, err := s.uc.UploadFile(ctx, "")
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler Upload file failed", err))
	}
	return responses.Response(ctx, http.StatusOK, attachments)
}

// UploadFile implements uploads.UseCase.
func (s *UploadHandler) AvatarUploadFile(ctx echo.Context) error {
	slog.Info("POST: AvatarUploadFile")
	attachments, err := s.uc.UploadFile(ctx, "avatar")
	if err != nil {
		return responses.ErrorResponse(ctx, http.StatusNotFound, responses.ErrorString("Handler Upload file failed", err))
	}
	return responses.Response(ctx, http.StatusOK, attachments)
}
func checkAttachmentTypeByCondition(attachableType string, condition []string) bool {
	temp := []string{constant.ATTACHMENT_POST, constant.ATTACHMENT_AVATAR, constant.ATTACHMENT_PROFILE}
	if condition == nil || len(condition) != 0 {
		condition = temp
	}
	return slices.Contains(condition, attachableType)
}
