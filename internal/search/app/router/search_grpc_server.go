package router

import (
	"context"
	"errors"

	v1 "github.com/dinhcanh303/go-microservices/api/search/v1"
	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/internal/search/usecases/searches"
	"github.com/google/wire"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type searchGRPCServer struct {
	v1.UnimplementedSearchServiceServer
	uc searches.UseCase
}

var _ v1.SearchServiceServer = (*searchGRPCServer)(nil)

var SearchGRPCServerSet = wire.NewSet(NewSearchGRPCServer)

func NewSearchGRPCServer(
	grpcServer *grpc.Server,
	uc searches.UseCase,
) v1.SearchServiceServer {
	svc := searchGRPCServer{
		uc: uc,
	}
	v1.RegisterSearchServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (s *searchGRPCServer) Search(ctx context.Context, request *v1.SearchRequest) (*v1.SearchResponse, error) {
	searchText := request.Q
	if searchText == "" {
		return nil, errors.New("key word search empty")
	}
	results, _ := s.uc.Search(searchText)
	return &v1.SearchResponse{
		Searches: lo.Map(results, func(item *domain.Search, _ int) *v1.Search {
			return &v1.Search{
				Id:         item.ID.String(),
				Name:       item.Name,
				Email:      item.Email,
				AvatarUrl:  item.AvatarUrl,
				ProfileUrl: item.ProfileUrl,
				FullName:   item.FullName,
				NickName:   item.NickName,
				Position:   item.Position,
			}
		}),
	}, nil
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
