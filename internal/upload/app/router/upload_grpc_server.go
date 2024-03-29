package router

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/upload/v1"
	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type uploadGRPCServer struct {
	v1.UnimplementedUploadServiceServer
	cfg *config.Config
	uc  uploads.UseCaseGRPC
}

var _ v1.UploadServiceServer = (*uploadGRPCServer)(nil)

var UploadServiceServer = wire.NewSet(NewGRPCUploadServer)

func NewGRPCUploadServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc uploads.UseCaseGRPC) v1.UploadServiceServer {
	svc := uploadGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	v1.RegisterUploadServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *uploadGRPCServer) GetAttachmentsByType(
	ctx context.Context,
	request *v1.GetAttachmentsByTypeRequest,
) (*v1.GetAttachmentsByTypeResponse, error) {
	slog.Info("GET: GetAttachmentsByType")
	attachableId, err := uuid.Parse(request.AttachableId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse")
	}
	attachments, err := g.uc.GetAttachmentsByType(ctx, request.AttachableType, attachableId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadGRPCServer.GetAttachmentsByType failed")
	}
	return &v1.GetAttachmentsByTypeResponse{
		Attachments: lo.Map(attachments, func(item *domain.Attachment, _ int) *v1.Attachment {
			return &v1.Attachment{
				Id:             item.ID.String(),
				AttachableType: item.AttachableType,
				AttachableId:   item.AttachableID.String(),
				UserId:         item.UserID.String(),
				EntityUploadId: item.EntityUploadID,
				Filename:       item.FileName,
				Extension:      item.Extension,
				MimeType:       item.MimeType,
				Folder:         item.Folder,
				Url:            item.URL,
				UrlThumbnail:   item.URLThumbnail,
				CreatedAt:      timestamppb.New(item.CreatedAt),
				UpdatedAt:      timestamppb.New(item.UpdatedAt),
			}
		}),
	}, nil
}

func (g *uploadGRPCServer) GetAttachmentsByOptional(
	ctx context.Context,
	request *v1.GetAttachmentsByOptionalRequest,
) (*v1.GetAttachmentsByOptionalResponse, error) {
	slog.Info("GET: GetAttachmentsByOptional")
	var err error
	var userId uuid.UUID
	if request.UserId != "" {
		userId, err = uuid.Parse(request.UserId)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse user id")
		}
	}
	attachment := &domain.Attachment{
		AttachableType: request.AttachableType,
		MimeType:       request.MimeType,
		UserID:         userId,
		EntityUploadID: request.EntityUploadId,
	}
	attachments, err := g.uc.GetAttachmentsByOptional(ctx, attachment, request.Limit, request.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "uploadGRPCServer.GetAttachmentsByOptional failed")
	}
	return &v1.GetAttachmentsByOptionalResponse{
		Attachments: lo.Map(attachments, func(item *domain.Attachment, _ int) *v1.Attachment {
			return &v1.Attachment{
				Id:             item.ID.String(),
				AttachableType: item.AttachableType,
				AttachableId:   item.AttachableID.String(),
				UserId:         item.UserID.String(),
				EntityUploadId: item.EntityUploadID,
				Filename:       item.FileName,
				Extension:      item.Extension,
				MimeType:       item.MimeType,
				Folder:         item.Folder,
				Url:            item.URL,
				UrlThumbnail:   item.URLThumbnail,
				CreatedAt:      timestamppb.New(item.CreatedAt),
				UpdatedAt:      timestamppb.New(item.UpdatedAt),
			}
		}),
	}, nil
}

func (g *uploadGRPCServer) GetAvatarUser(
	ctx context.Context,
	request *v1.GetAvatarUserRequest,
) (*v1.GetAvatarUserResponse, error) {
	slog.Info("GET: GetAvatarUser")
	attachableId, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse")
	}
	attachment, err := g.uc.GetLastAttachmentByType(ctx, constant.ATTACHMENT_AVATAR, attachableId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadGRPCServer.GetAvatarUser failed")
	}

	return &v1.GetAvatarUserResponse{
		Attachment: &v1.Attachment{
			Id:             attachment.ID.String(),
			AttachableType: attachment.AttachableType,
			AttachableId:   attachment.AttachableID.String(),
			UserId:         attachment.UserID.String(),
			EntityUploadId: attachment.EntityUploadID,
			Filename:       attachment.FileName,
			Extension:      attachment.Extension,
			MimeType:       attachment.MimeType,
			Folder:         attachment.Folder,
			Url:            attachment.URL,
			UrlThumbnail:   attachment.URLThumbnail,
			CreatedAt:      timestamppb.New(attachment.CreatedAt),
			UpdatedAt:      timestamppb.New(attachment.UpdatedAt),
		},
	}, nil
}
