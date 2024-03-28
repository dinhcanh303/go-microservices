package grpc

import (
	"context"
	"log/slog"

	v1 "github.com/dinhcanh303/go-microservices/api/upload/v1"
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type uploadGRPCClient struct {
	conn *grpc.ClientConn
}

// GetAvatarUser implements domain.UploadDomainService.
func (u *uploadGRPCClient) GetAvatarUser(ctx context.Context, userId uuid.UUID) (*domainUpload.Attachment, error) {
	client := v1.NewUploadServiceClient(u.conn)
	res, err := client.GetAvatarUser(ctx, &v1.GetAvatarUserRequest{
		Id: userId.String(),
	})
	if err != nil {
		slog.Warn("uploadGRPCClient.GetAvatarUser failed", err)
		return &domainUpload.Attachment{}, nil
	}
	return &domainUpload.Attachment{
		ID:             uuid.MustParse(res.Attachment.Id),
		UserID:         uuid.MustParse(res.Attachment.UserId),
		AttachableType: res.Attachment.AttachableType,
		AttachableID:   uuid.MustParse(res.Attachment.AttachableId),
		FileName:       res.Attachment.Filename,
		Extension:      res.Attachment.Extension,
		MimeType:       res.Attachment.MimeType,
		Folder:         res.Attachment.Folder,
		URL:            res.Attachment.Url,
		URLThumbnail:   res.Attachment.UrlThumbnail,
		CreatedAt:      res.Attachment.CreatedAt.AsTime(),
		UpdatedAt:      res.Attachment.UpdatedAt.AsTime(),
	}, nil
}

var UploadGRPCClientSet = wire.NewSet(NewGRPCUploadClient)

var _ domain.UploadDomainService = (*uploadGRPCClient)(nil)

func NewGRPCUploadClient(cfg *config.Config) (domain.UploadDomainService, error) {
	conn, err := grpc.Dial(cfg.UploadClient.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &uploadGRPCClient{
		conn: conn,
	}, nil
}
