package router

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/search/domain"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/elastic"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type searchGRPCServer struct {
	gen.UnimplementedSearchServiceServer
	elastic elastic.ElasticSearch
}

var _ gen.SearchServiceServer = (*searchGRPCServer)(nil)

var SearchGRPCServerSet = wire.NewSet(NewSearchGRPCServer)

func NewSearchGRPCServer(
	grpcServer *grpc.Server,
	elastic elastic.ElasticSearch,
) gen.SearchServiceServer {
	svc := searchGRPCServer{
		elastic: elastic,
	}
	gen.RegisterSearchServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}
func (c *searchGRPCServer) Search(ctx context.Context, request *gen.SearchRequest) (*gen.SearchResponse, error) {
	keyWord := request.KeyWord
	slog.Info("key word", keyWord)
	if keyWord == "" {
		return nil, errors.New("key word search empty")
	}
	results, err := searchElastic(c.elastic, constant.ELASTIC_SEARCH_INDEX, "group_name", keyWord)
	if err != nil {
		return nil, err
	}
	slog.Info("value search", results)
	return &gen.SearchResponse{
		Search: "Hehe",
	}, nil
}

func searchElastic(es elastic.ElasticSearch, indexName, fieldName, keyWord string) ([]domain.GroupSearch, error) {
	// var users []domain.UserSearch
	var groups []domain.GroupSearch

	var buf strings.Builder
	buf.WriteString(`{
		"query": {
			"match": {
				"` + fieldName + `": "` + keyWord + `"
			}
		}
	}`)
	results, err := es.Search(indexName, buf.String())
	if err != nil {
		return nil, err
	}
	slog.Info("DATA_ALL::", results)
	// Loop Response
	// Map And Append Document
	if results["hits"] == nil {
		return nil, nil
	}
	slog.Info("DATA::", results["hits"])

	for _, hit := range results["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var group = domain.GroupSearch{}

		utils.Mapping(hit.(map[string]interface{})["_source"], &group)
		groups = append(groups, group)
	}
	return groups, nil

}
