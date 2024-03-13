package router

import (
	"context"
	"time"

	"github.com/dinhcanh303/go-microservices/internal/notification/domain"
	"github.com/dinhcanh303/go-microservices/internal/notification/usecases/notifications"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type notiGRPCServer struct {
	gen.UnimplementedNotiServiceServer
	uc            notifications.UseCase
	authDomainSvc domain.AuthDomainService
}

var _ gen.NotiServiceServer = (*notiGRPCServer)(nil)

var NotiGRPCServerSet = wire.NewSet(NewNotiGRPCServer)

func NewNotiGRPCServer(
	grpcServer *grpc.Server,
	uc notifications.UseCase,
	authDomainSvc domain.AuthDomainService,
) gen.NotiServiceServer {
	svc := notiGRPCServer{
		uc:            uc,
		authDomainSvc: authDomainSvc,
	}
	gen.RegisterNotiServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (n *notiGRPCServer) GetNotifications(ctx context.Context, request *gen.GetNotificationsRequest) (*gen.GetNotificationsResponse, error) {
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	notifications, err := n.uc.GetNotificationsByUserId(ctx, user.ID, int(request.Limit), int(request.Offset), sharedkernel.GetNotiOptions{
		Unread: request.Options,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetNotificationsByUserId failed")
	}
	return &gen.GetNotificationsResponse{
		Notifications: lo.Map(notifications, func(item domain.Notification, _ int) *gen.Notification {
			str, _ := utils.MapToProtobufStruct(item.Data)
			var readAt *timestamppb.Timestamp
			if item.ReadAt != nil {
				readAt = timestamppb.New(item.CreatedAt)
			}
			actorProfile, _ := n.authDomainSvc.GetProfile(ctx, item.ActorID)
			return &gen.Notification{
				Id:         item.ID,
				ActorId:    item.ActorID,
				SenderId:   item.SenderID,
				Actor:      actorProfile.User,
				Data:       str,
				ObjectType: item.ObjectType,
				ObjectId:   item.ObjectID,
				ReadAt:     readAt,
				CreatedAt:  timestamppb.New(item.CreatedAt),
				UpdatedAt:  timestamppb.New(item.UpdatedAt),
			}
		}),
	}, err
}

func (n *notiGRPCServer) ReadNotification(ctx context.Context, request *gen.ReadNotificationRequest) (*gen.ReadNotificationResponse, error) {
	_, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	timeNow := time.Now()
	err = n.uc.ReadNotification(ctx, &domain.Notification{
		ID:     request.Id,
		ReadAt: &timeNow,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uc.ReadNotification failed")
	}
	return nil, nil
}
