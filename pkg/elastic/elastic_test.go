package elastic

import (
	"testing"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestElasticsearchConnect(t *testing.T) {
	elastic, err := ConnectElasticsearch()
	require.NoError(t, err)
	require.NotEmpty(t, elastic)
}

func TestInsertItemElasticsearch(t *testing.T) {

}

func TestRemoveItemElasticsearch(t *testing.T) {

}

func TestPingElasticsearch(t *testing.T) {

}

func ConnectElasticsearch() (ElasticSearch, error) {
	err := utils.LoadFileEnvOnLocal()
	if err != nil {
		return nil, err
	}
	config, err := configs.NewConfigElasticSearch()
	if err != nil {
		return nil, err
	}
	elastic, err := NewElasticSearch(*config)
	if err != nil {
		return nil, err
	}
	return elastic, nil
}
