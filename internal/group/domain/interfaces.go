package domain

import "context"

type GroupRepo interface {
	Get(ctx context.Context, uuid string) (*Group, error)
	GetWithUnscoped(ctx context.Context, uuid string) (*Group, error)
	Create(ctx context.Context, group *Group) (*Group, error)
	Update(ctx context.Context, group *Group) (*Group, error)
	Delete(ctx context.Context, uuid string) (bool, error)
}
type AttachmentDomainService interface {
}
