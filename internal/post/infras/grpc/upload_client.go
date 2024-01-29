package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type uploadGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.UploadDomainService = (*uploadGRPCClient)(nil)

var UploadGRPCClientSet = wire.NewSet(NewGRPCUploadClient)

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

// GetAttachmentsByType implements domain.UploadDomainService.
func (u *uploadGRPCClient) GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error) {
	client := gen.NewUploadServiceClient(u.conn)
	res, err := client.GetAttachmentsByType(ctx, &gen.GetAttachmentsByTypeRequest{
		AttachableType: attachableType,
		AttachableId:   attachableId.String(),
	})
	results := make([]*domainUpload.Attachment, 0)
	if err != nil {
		slog.Warn("uploadGRPCClient.GetAttachmentsByType failed", err)
		return results, nil
	}
	for _, item := range res.Attachments {
		results = append(results, &domainUpload.Attachment{
			ID:             uuid.MustParse(item.Id),
			UserID:         uuid.MustParse(item.UserId),
			AttachableType: item.AttachableType,
			AttachableID:   uuid.MustParse(item.AttachableId),
			FileName:       item.Filename,
			Extension:      item.Extension,
			MimeType:       item.MimeType,
			Folder:         item.Folder,
			URL:            item.Url,
			URLThumbnail:   item.UrlThumbnail,
			CreatedAt:      item.CreatedAt.AsTime(),
			UpdatedAt:      item.UpdatedAt.AsTime(),
		})
	}
	return results, nil
}
