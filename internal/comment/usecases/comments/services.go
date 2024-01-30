package comments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type service struct {
	commentRepo     CommentRepo
	redis           redis.RedisEngine
	likeDomainSvc   domain.LikeDomainService
	uploadDomainSvc domain.UploadDomainService
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(commentRepo CommentRepo,
	redis redis.RedisEngine,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService) UseCase {
	return &service{
		commentRepo:     commentRepo,
		redis:           redis,
		likeDomainSvc:   likeDomainSvc,
		uploadDomainSvc: uploadDomainSvc,
	}
}

// GetCommentsByCommentID implements UseCase.
func (s *service) GetCommentsByCommentID(ctx context.Context,
	commentId uuid.UUID,
	userId uuid.UUID,
	limit int32, offset int32) ([]*domain.CommentHasMetadata, error) {
	var comments []*domain.Comment
	keyCache := constant.CacheCommentsByCommentId + commentId.String() +
		constant.CacheLimit + utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(comments, s.redis, keyCache)
	if err != nil {
		comments, err = s.commentRepo.GetCommentsByCommentID(ctx, commentId, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetCommentsByCommentID")
		}
		err = s.redis.Set(keyCache, comments)
		if err != nil {
			slog.Error("set cached comments failed", err)
		}
	}
	var commentsHasMetadata []*domain.CommentHasMetadata
	for _, comment := range comments {
		likeInfo, err := s.likeDomainSvc.GetLikesInfoByCommentID(ctx, comment.ID, userId)
		if err != nil {
			likeInfo = nil
		}
		attachments, err := s.uploadDomainSvc.GetAttachmentsByType(ctx, "Comment", comment.ID)
		if err != nil {
			attachments = make([]*domainUpload.Attachment, 0)
		}
		commentsHasMetadata = append(commentsHasMetadata, &domain.CommentHasMetadata{
			ID:              comment.ID,
			UserID:          comment.UserID,
			ReplyID:         comment.ReplyID,
			Content:         comment.Content,
			PostID:          comment.PostID,
			ParentCommentID: comment.ParentCommentID,
			Likes:           likeInfo,
			TagIDs:          comment.TagIDs,
			Attachments:     attachments,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		})
	}
	return commentsHasMetadata, nil
}

// CreateComment implements UseCase.
func (s *service) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	comment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateComment")
	}
	//Invalidate cache
	err = s.redis.InvalidatePrefix(constant.CacheComments)
	if err != nil {
		slog.Error("Invalidate cache comments failed", err)
	}
	return comment, nil
}

// DeleteComment implements UseCase.
func (s *service) DeleteComment(ctx context.Context, id uuid.UUID) (bool, error) {
	isDelete, err := s.commentRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteComment")
	}
	//Invalidate cache
	err = s.redis.InvalidatePrefix(constant.CacheComments)
	if err != nil {
		slog.Error("Invalidate cache comments failed", err)
	}
	return isDelete, nil
}

// DeleteAllCommentByPostID implements UseCase.
func (s *service) DeleteAllCommentByPostID(ctx context.Context, postId uuid.UUID) (bool, error) {
	isDelete, err := s.commentRepo.DeleteAllByPostID(ctx, postId)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteAllCommentByPostID")
	}
	//Invalidate cache
	err = s.redis.InvalidatePrefix(constant.CacheComments)
	if err != nil {
		slog.Error("Invalidate cache comments failed", err)
	}
	return isDelete, nil
}

// GetComment implements UseCase.
func (s *service) GetComment(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	comment, err := s.commentRepo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetComment")
	}
	return comment, nil
}

// GetCommentsByPostID implements UseCase.
func (s *service) GetCommentsByPostID(ctx context.Context, postId, userId uuid.UUID, limit, offset int32) ([]*sharedkernel.CommentHasChildren, error) {
	var comments []*domain.Comment
	keyCache := constant.CacheCommentsByPostId + postId.String() +
		constant.CacheLimit + utils.String(limit) + constant.CacheOffset + utils.String(offset)
	err := utils.HandleHitCache(comments, s.redis, keyCache)
	if err != nil {
		comments, err = s.commentRepo.GetCommentsByPostID(ctx, postId, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "service.GetCommentsByPostID")
		}
		err = s.redis.Set(keyCache, comments)
		if err != nil {
			slog.Error("set cached comments failed", err)
		}
	}

	commentMap := make(map[uuid.UUID]*sharedkernel.CommentHasChildren)
	var commentsHasChildren []*sharedkernel.CommentHasChildren
	slog.Info("COMMENT::", comments)
	for _, comment := range comments {
		likeInfo, err := s.likeDomainSvc.GetLikesInfoByCommentID(ctx, comment.ID, userId)
		slog.Info("LIKE INFO::", likeInfo)
		if err != nil {
			likeInfo = nil
		}
		attachments, err := s.uploadDomainSvc.GetAttachmentsByType(ctx, "Comment", comment.ID)
		if err != nil {
			attachments = make([]*domainUpload.Attachment, 0)
		}
		commentHasChildren := &sharedkernel.CommentHasChildren{
			ID:              comment.ID,
			UserID:          comment.UserID,
			ReplyID:         comment.ReplyID,
			TagIDs:          comment.TagIDs,
			Content:         comment.Content,
			PostID:          comment.PostID,
			ParentCommentID: comment.ParentCommentID,
			Likes:           likeInfo,
			Attachments:     attachments,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		}
		results := &domain.CommentHasMetadata{
			ID:              comment.ID,
			UserID:          comment.UserID,
			ReplyID:         comment.ReplyID,
			TagIDs:          comment.TagIDs,
			Content:         comment.Content,
			PostID:          comment.PostID,
			ParentCommentID: comment.ParentCommentID,
			Likes:           likeInfo,
			Attachments:     attachments,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		}
		commentMap[comment.ID] = commentHasChildren
		if comment.ParentCommentID.UUID.String() != constant.NullUUID {
			parentComment, exists := commentMap[comment.ParentCommentID.UUID]
			if exists {
				parentComment.Children = append(parentComment.Children, results)
			}
		} else {
			commentsHasChildren = append(commentsHasChildren, commentHasChildren)
		}
	}
	return commentsHasChildren, nil
}

// UpdateComment implements UseCase.
func (s *service) UpdateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	comment, err := s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, errors.Wrap(err, "service.UpdateComment")
	}
	//Invalidate cache
	err = s.redis.InvalidatePrefix(constant.CacheComments)
	if err != nil {
		slog.Error("Invalidate cache comments failed", err)
	}

	return comment, nil
}

// CountCommentByCommentID implements UseCase.
func (s *service) CountCommentByCommentID(ctx context.Context, commentId uuid.UUID) (int64, error) {
	count, err := s.commentRepo.CountByCommentID(ctx, commentId)
	if err != nil {
		return 0, errors.Wrap(err, "service.UpdateComment")
	}
	return count, nil
}

// CountCommentByPostID implements UseCase.
func (s *service) CountCommentByPostID(ctx context.Context, postId uuid.UUID) (int64, error) {
	var count int64
	keyCache := constant.CacheCommentsCountByPostId + postId.String()
	err := utils.HandleHitCache(count, s.redis, keyCache)
	if err != nil {
		count, err = s.commentRepo.CountByPostID(ctx, postId)
		if err != nil {
			return 0, errors.Wrap(err, "service.UpdateComment")
		}
		err = s.redis.Set(keyCache, count)
		if err != nil {
			slog.Error("set cache count comment by post failed", err)
		}
	}

	return count, nil
}
