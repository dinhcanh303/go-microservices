package router

import (
	"context"
	"errors"

	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/meili"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type searchGRPCServer struct {
	gen.UnimplementedSearchServiceServer
	meili meili.MeiliSearch
}

var _ gen.SearchServiceServer = (*searchGRPCServer)(nil)

var SearchGRPCServerSet = wire.NewSet(NewSearchGRPCServer)

func NewSearchGRPCServer(
	grpcServer *grpc.Server,
	meili meili.MeiliSearch,
) gen.SearchServiceServer {
	svc := searchGRPCServer{
		meili: meili,
	}
	gen.RegisterSearchServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (c *searchGRPCServer) Search(ctx context.Context, request *gen.SearchRequest) (*gen.SearchResponse, error) {
	searchText := request.Q
	if searchText == "" {
		return nil, errors.New("key word search empty")
	}
	hits, err := meiliSearch(c.meili, constant.MEILI_SEARCH_INDEX, searchText, []string{"name", "email"})
	if err != nil {
		return nil, err
	}
	if hits == nil {
		return nil, nil
	}
	var results []searchRes
	for _, hit := range hits {
		var res = searchRes{}
		utils.Mapping(hit, &res)
		results = append(results, res)
	}
	return &gen.SearchResponse{
		Search: lo.Map(results, func(item searchRes, _ int) *gen.Search {
			return &gen.Search{
				Id:        item.ID.String(),
				Name:      item.Name,
				Email:     item.Email,
				AvatarUrl: item.Email,
				Type:      item.Type,
			}
		}),
	}, nil
}

type searchRes struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
	Type   string    `json:"type"`
}

func meiliSearch(ml meili.MeiliSearch, indexName, keyWord string, fieldName []string) ([]interface{}, error) {
	task, err := ml.Search(indexName, keyWord, &meilisearch.SearchRequest{
		Limit:                20,
		AttributesToSearchOn: fieldName,
	})
	if err != nil {
		return nil, err
	}
	return task.Hits, nil
}

// func searchElastic(es elastic.ElasticSearch, indexName, fieldName, keyWord string) ([]domain.GroupSearch, error) {
// 	// var users []domain.UserSearch
// 	var groups []domain.GroupSearch

// 	var buf strings.Builder
// 	buf.WriteString(`{
// 		"query": {
// 			"match": {
// 				"` + fieldName + `": "` + keyWord + `"
// 			}
// 		}
// 	}`)
// 	results, err := es.Search(indexName, buf.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	slog.Info("DATA_ALL::", results)
// 	// Loop Response
// 	// Map And Append Document
// 	if results["hits"] == nil {
// 		return nil, nil
// 	}
// 	slog.Info("DATA::", results["hits"])

// 	for _, hit := range results["hits"].(map[string]interface{})["hits"].([]interface{}) {
// 		var group = domain.GroupSearch{}

// 		utils.Mapping(hit.(map[string]interface{})["_source"], &group)
// 		groups = append(groups, group)
// 	}
// 	return groups, nil

// }
