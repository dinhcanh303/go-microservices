package posts

import (
	"context"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/feed/domain"
	domain2 "github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type usecase struct {
	postDomainSvc  domain.PostDomainService
	groupDomainSvc domain.GroupDomainService
}

// NewFeed implements UseCase.
func (p *usecase) NewFeed(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain2.PostExtra, error) {
	groupIds, err := p.groupDomainSvc.GetAllGroupIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "groupDomainSvc.GetAllGroupIdByUserId")
	}
	posts, err := p.postDomainSvc.GetPostsByFeed(ctx, "", strings.Join(groupIds, ","), limit, offset)
	panic("unimplemented")
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(postDomainSvc domain.PostDomainService) UseCase {
	return &usecase{
		postDomainSvc: postDomainSvc,
	}
}
