package router

import (
	"context"
	"errors"

	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/internal/search/usecases/searches"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type searchGRPCServer struct {
	gen.UnimplementedSearchServiceServer
	uc searches.UseCase
}

var _ gen.SearchServiceServer = (*searchGRPCServer)(nil)

var SearchGRPCServerSet = wire.NewSet(NewSearchGRPCServer)

func NewSearchGRPCServer(
	grpcServer *grpc.Server,
	uc searches.UseCase,
) gen.SearchServiceServer {
	svc := searchGRPCServer{
		uc: uc,
	}
	gen.RegisterSearchServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (s *searchGRPCServer) Search(ctx context.Context, request *gen.SearchRequest) (*gen.SearchResponse, error) {
	searchText := request.Q
	if searchText == "" {
		return nil, errors.New("key word search empty")
	}
	results, _ := s.uc.Search(searchText)
	return &gen.SearchResponse{
		Searches: lo.Map(results, func(item *domain.Search, _ int) *gen.Search {
			return &gen.Search{
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
