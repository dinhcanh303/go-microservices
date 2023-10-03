package domain

import "context"

type GroupRepo interface {
	Get(ctx context.Context, uuid string) (*Attachment, error)
	GetAll() ([]*Attachment, error)
	GetWithUnscoped(ctx context.Context, uuid string) (*Attachment, error)
	Create(ctx context.Context, group *Attachment) (*Attachment, error)
	Update(ctx context.Context, group *Attachment) (Attachment, error)
	Delete(ctx context.Context, uuid string) (bool, error)
}
