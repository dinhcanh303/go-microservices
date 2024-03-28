package sharedkernel

import (
	v1l "github.com/dinhcanh303/go-microservices/api/like/v1"
	v1u "github.com/dinhcanh303/go-microservices/api/upload/v1"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EntityAttachmentToProtobuf(attachments []*domainUpload.Attachment) []*v1u.Attachment {
	return lo.Map(attachments, func(attachment *domainUpload.Attachment, _ int) *v1u.Attachment {
		return &v1u.Attachment{
			Id:             attachment.ID.String(),
			UserId:         attachment.UserID.String(),
			AttachableType: attachment.AttachableType,
			AttachableId:   attachment.AttachableID.String(),
			Filename:       attachment.FileName,
			Url:            attachment.URL,
			UrlThumbnail:   attachment.URLThumbnail,
			Extension:      attachment.Extension,
			MimeType:       attachment.MimeType,
			Folder:         attachment.Folder,
			CreatedAt:      timestamppb.New(attachment.CreatedAt),
			UpdatedAt:      timestamppb.New(attachment.UpdatedAt),
		}
	})
}
func EntityLikeToProtobuf(like *domainLike.LikesInfo) *v1l.LikeInfo {
	return &v1l.LikeInfo{
		YourLikedEmoji:    like.YourLikedEmoji,
		YourLike:          like.YourLike,
		OthersLikedEmojis: like.OthersLikedEmojis,
		OthersLikes:       like.OthersLikes,
	}
}
