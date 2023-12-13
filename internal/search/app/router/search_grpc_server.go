package router

import (
	"context"
	"errors"

	"github.com/dinhcanh303/go-microservices/pkg/elastic"
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
	if keyWord == "" {
		return nil, errors.New("Key Word search empty")
	}
	return nil, nil
}

// func searchElastic(es elastic.ElasticSearch, indexName, fieldName, keyWord string) {
// 	// var users []domain.UserSearch
// 	// var groups []domain.GroupSearch

// 	var buf strings.Builder
// 	buf.WriteString(`{
// 		"query": {
// 			"bool": {
// 				"should": [
// 				{
// 					"match": {
// 					"user_name": "canh"
// 					}
// 				},
// 				{
// 					"match": {
// 					"group_name": "canh"
// 					}
// 				}
// 				]
// 			}
// 		}
// 	}`)
// 	// results, err := es.Search(indexName, buf.String())
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// }
