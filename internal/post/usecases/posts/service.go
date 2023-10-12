package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type usecase struct {
	postRepo PostRepo
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(postRepo PostRepo) UseCase {
	return &usecase{
		postRepo: postRepo,
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
