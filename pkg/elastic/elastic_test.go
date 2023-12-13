package elastic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	URL      = "http://192.168.100.86:9200/"
	UserName = "elastic"
	Password = "elastic"
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
	esConn := &ElasticSearchConn{
		Url:      URL,
		UserName: UserName,
		Password: Password,
	}
	elastic, err := NewElasticSearch(*esConn)
	if err != nil {
		return nil, err
	}
	return elastic, nil
}
