package searches

import (
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
)

type UseCase interface {
	Search(search string) ([]*domain.Search, error)
}
