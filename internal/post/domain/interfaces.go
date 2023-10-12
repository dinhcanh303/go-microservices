package domain

import "context"

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context)
	}
)
