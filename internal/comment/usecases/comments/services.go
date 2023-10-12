package comments

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/comment/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type service struct {
	commentRepo CommentRepo
}

var _ UseCase = (*service)(nil)

func NewService(commentRepo CommentRepo) UseCase {
	return &service{
		commentRepo: commentRepo,
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
func (s *service) GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Comment, error) {
	comments, err := s.commentRepo.GetCommentByPostID(ctx, postId)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetCommentsByPostID")
	}
	return comments, nil
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
