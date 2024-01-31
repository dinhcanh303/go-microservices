package searches

import (
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/meili"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/wire"
	"github.com/meilisearch/meilisearch-go"
)

type service struct {
	meili meili.MeiliSearch
}

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(NewService)

func NewService(meili meili.MeiliSearch) UseCase {
	return &service{
		meili: meili,
	}
}

// Search implements UseCase.
func (s *service) Search(search string) ([]*domain.Search, error) {
	var results = make([]*domain.Search, 0)
	searchUserRes, _ := meiliSearch(s.meili, constant.MeiliSearchDBUserIndex,
		search, []string{"name", "email", "phone"})
	searchGroupRes, _ := meiliSearch(s.meili, constant.MeiliSearchDBGroupIndex,
		search, []string{"name"})
	results = append(searchUserRes, searchGroupRes...)
	return results, nil
}
func meiliSearch(ml meili.MeiliSearch, indexName,
	keyWord string, fieldName []string) ([]*domain.Search, error) {
	task, err := ml.Search(indexName, keyWord, &meilisearch.SearchRequest{
		Limit:                10,
		AttributesToSearchOn: fieldName,
	})
	if err != nil {
		return nil, err
	}

	if task.Hits == nil {
		return nil, nil
	}
	var results []*domain.Search
	for _, hit := range task.Hits {
		var res = domain.Search{}
		utils.Mapping(hit, &res)
		results = append(results, &res)
	}
	return results, nil
}
