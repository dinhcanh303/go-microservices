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

// GetPostsByFeed implements UseCase.
func (*usecase) GetPostsByFeed(ctx context.Context, userIds, groupIds string, limit int32, offset int32) ([]*domain.PostExtra, error) {
	panic("unimplemented")
}

// GetPostsByGroupId implements UseCase.
func (*usecase) GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit int32, offset int32) ([]*domain.PostExtra, error) {
	panic("unimplemented")
}

// GetPostsByUserId implements UseCase.
func (*usecase) GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit int32, offset int32) ([]*domain.PostExtra, error) {

	panic("unimplemented")
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(postRepo PostRepo,
) UseCase {
	return &usecase{
		postRepo: postRepo,
	}
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

// UpdatePost implements UseCase.
func (uc *usecase) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	post, err := uc.postRepo.Update(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "postRepo.UpdatePost")
	}
	return post, nil
}
