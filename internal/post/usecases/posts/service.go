package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type usecase struct {
	postRepo         PostRepo
	commentDomainSvc domain.CommentDomainService
	likeDomainSvc    domain.LikeDomainService
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(postRepo PostRepo,
	commentDomainSvc domain.CommentDomainService,
	likeDomainSvc domain.LikeDomainService) UseCase {
	return &usecase{
		postRepo:         postRepo,
		commentDomainSvc: commentDomainSvc,
		likeDomainSvc:    likeDomainSvc,
	}
}

// // Count implements UseCase.
//
//	func (uc *usecase) Count(ctx context.Context) (uint64, error) {
//		count, err := uc.postRepo.Count(ctx)
//		if err != nil {
//			return 0, errors.Wrap(err, "postRepo.Count")
//		}
//		return count, nil
//	}
//

// GetPostExtra implements UseCase.
func (uc *usecase) GetPostExtra(ctx context.Context, id uuid.UUID) (*domain.PostExtra, error) {
	post, err := uc.GetPost(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Get")
	}
	comments, err := uc.commentDomainSvc.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		comments = nil
	}
	likes, err := uc.likeDomainSvc.GetLikesByPostID(ctx, post.ID)
	if err != nil {
		likes = nil
	}
	return &domain.PostExtra{
		ID:        post.ID,
		Status:    post.Status,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		GroupID:   post.GroupID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Comments:  comments,
		Likes:     likes,
	}, nil
}

// GetByGroupId implements UseCase.
func (uc *usecase) GetByGroupId(ctx context.Context, groupId uuid.UUID) ([]*domain.Post, error) {
	posts, err := uc.postRepo.GetByGroupId(ctx, groupId)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.GetByGroupId")
	}
	return posts, nil
}

// GetByUserId implements UseCase.
func (uc *usecase) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*domain.Post, error) {
	posts, err := uc.postRepo.GetByGroupId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.GetByUserId")
	}
	return posts, nil
}

// CreatePost implements UseCase.
func (uc *usecase) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Create(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Create")
	}
	return post, nil
}

// DeletePost implements UseCase.
func (uc *usecase) DeletePost(ctx context.Context, id uuid.UUID) (bool, error) {
	isDelete, err := uc.postRepo.Delete(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "postRepo.Delete")
	}
	return isDelete, nil
}

// GetPost implements UseCase.
func (uc *usecase) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	post, err := uc.postRepo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.Get")
	}
	return post, nil
}

// ListPost implements UseCase.
func (uc *usecase) ListPost(ctx context.Context, offset int, limit int) ([]*domain.Post, error) {
	posts, err := uc.postRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.List")
	}
	return posts, nil
}

// UpdatePost implements UseCase.
func (uc *usecase) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Update(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.GetWithUnscoped")
	}
	return post, nil
}
