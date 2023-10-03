package domain

import "context"

type PostRepo interface {
	Get(ctx context.Context, uuid string) (*Post, error)
	GetWithUnscoped(ctx context.Context, uuid string) (*Post, error)
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (Post, error)
	UpdateWithUnscoped(ctx context.Context, post *Post) (Post, error)
	Delete(ctx context.Context, uuid string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*Post, error)
	Count(ctx context.Context) (uint64, error)
}
