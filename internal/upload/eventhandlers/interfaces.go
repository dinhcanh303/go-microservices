package eventhandlers

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
)

type (
	GroupUploadEventHandler interface {
		Handle(context.Context, event.GroupUploadBegin)
	}
	PostUploadEventHandler interface {
		Handle(context.Context, event.PostUploadBegin)
	}
	CommentUploadEventHandler interface {
		Handle(context.Context, event.CommentUploadBegin)
	}
)
