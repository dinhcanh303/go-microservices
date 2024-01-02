package router

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/upload/config"
	"github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/internal/upload/usecases/uploads"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type uploadGRPCServer struct {
	gen.UnimplementedUploadServiceServer
	cfg *config.Config
	uc  uploads.UseCaseGRPC
}

var _ gen.UploadServiceServer = (*uploadGRPCServer)(nil)

var UploadServiceServer = wire.NewSet(NewGRPCUploadServer)

func NewGRPCUploadServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc uploads.UseCaseGRPC) gen.UploadServiceServer {
	svc := uploadGRPCServer{
		cfg: cfg,
		uc:  uc,
	}
	gen.RegisterUploadServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (g *uploadGRPCServer) GetAttachmentsByType(
	ctx context.Context,
	request *gen.GetAttachmentsByTypeRequest,
) (*gen.GetAttachmentsByTypeResponse, error) {
	slog.Info("GET: GetAttachmentsByType")
	attachableId, err := uuid.Parse(request.AttachableId)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse")
	}
	attachments, err := g.uc.GetAttachmentsByType(ctx, request.AttachableType, attachableId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadGRPCServer.GetAttachmentsByType failed")
	}
	return &gen.GetAttachmentsByTypeResponse{
		Attachments: lo.Map(attachments, func(item *domain.Attachment, _ int) *gen.Attachment {
			return &gen.Attachment{
				Id:             item.ID.String(),
				AttachableType: item.AttachableType,
				AttachableId:   item.AttachableID.String(),
				UserId:         item.UserID.String(),
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
	request *gen.GetAvatarUserRequest,
) (*gen.GetAvatarUserResponse, error) {
	slog.Info("GET: GetAvatarUser")
	attachableId, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse")
	}
	attachment, err := g.uc.GetLastAttachmentByType(ctx, "Attachment/Avatar", attachableId)
	if err != nil {
		return nil, errors.Wrap(err, "uploadGRPCServer.GetAvatarUser failed")
	}

	return &gen.GetAvatarUserResponse{
		Attachment: &gen.Attachment{
			Id:             attachment.ID.String(),
			AttachableType: attachment.AttachableType,
			AttachableId:   attachment.AttachableID.String(),
			UserId:         attachment.UserID.String(),
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
