package comments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

const UUID_NULL string = "00000000-0000-0000-0000-000000000000"

type service struct {
	commentRepo     CommentRepo
	likeDomainSvc   domain.LikeDomainService
	uploadDomainSvc domain.UploadDomainService
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(commentRepo CommentRepo,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService) UseCase {
	return &service{
		commentRepo:     commentRepo,
		likeDomainSvc:   likeDomainSvc,
		uploadDomainSvc: uploadDomainSvc,
	}
}

// CreateComment implements UseCase.
func (s *service) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	comment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, errors.Wrap(err, "service.CreateComment")
	}
	return comment, nil
}

// DeleteComment implements UseCase.
func (s *service) DeleteComment(ctx context.Context, id uuid.UUID) (bool, error) {
	isDelete, err := s.commentRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteComment")
	}
	return isDelete, nil
}

// DeleteAllCommentByPostID implements UseCase.
func (s *service) DeleteAllCommentByPostID(ctx context.Context, postId uuid.UUID) (bool, error) {
	isDelete, err := s.commentRepo.DeleteAllByPostID(ctx, postId)
	if err != nil {
		return false, errors.Wrap(err, "service.DeleteAllCommentByPostID")
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
func (s *service) GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error) {
	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetCommentsByPostID")
	}
	commentMap := make(map[uuid.UUID]*sharedkernel.CommentHasChildren)
	var commentsHasChildren []*sharedkernel.CommentHasChildren
	slog.Info("COMMENT::", comments)
	for i, comment := range comments {
		likes, err := s.likeDomainSvc.GetLikesByCommentID(ctx, comment.ID)
		slog.Info("LIKE::", likes)
		if err != nil {
			likes = make([]*domainLike.Like, 0)
		}
		attachments, err := s.uploadDomainSvc.GetAttachmentsByType(ctx, "Attachment/Comment", comment.ID)
		slog.Info("LIKE::", likes)
		if err != nil {
			attachments = make([]*domainUpload.Attachment, 0)
		}
		commentHasChildren := &sharedkernel.CommentHasChildren{
			ID:              comments[i].ID,
			UserID:          comments[i].UserID,
			ReplyToID:       comments[i].ReplyToID,
			Content:         comments[i].Content,
			PostID:          comments[i].PostID,
			ParentCommentID: comments[i].ParentCommentID,
			Likes:           likes,
			Attachments:     attachments,
			CreatedAt:       comments[i].CreatedAt,
			UpdatedAt:       comments[i].UpdatedAt,
		}
		results := &domain.CommentHasLike{
			ID:              comments[i].ID,
			UserID:          comments[i].UserID,
			ReplyToID:       comments[i].ReplyToID,
			Content:         comments[i].Content,
			PostID:          comments[i].PostID,
			ParentCommentID: comments[i].ParentCommentID,
			Likes:           likes,
			Attachments:     attachments,
			CreatedAt:       comments[i].CreatedAt,
			UpdatedAt:       comments[i].UpdatedAt,
		}
		commentMap[comment.ID] = commentHasChildren
		if comment.ParentCommentID.UUID.String() != UUID_NULL {
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
	count, err := s.commentRepo.CountByPostID(ctx, postId)
	if err != nil {
		return 0, errors.Wrap(err, "service.UpdateComment")
	}
	return count, nil
}
