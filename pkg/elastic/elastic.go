package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticSearchConn struct {
	Url      string
	UserName string
	Password string
}
type elasticSearch struct {
	Client *elasticsearch.Client
}

var _ ElasticSearch = (*elasticSearch)(nil)

func NewElasticSearch(elasticSearchConn ElasticSearchConn) (ElasticSearch, error) {
	cfgES := elasticsearch.Config{
		Addresses: []string{elasticSearchConn.Url},
		Username:  elasticSearchConn.UserName,
		Password:  elasticSearchConn.Password,
	}
	// Create New Elasticsearch Client
	esClient, err := elasticsearch.NewClient(cfgES)
	if err != nil {
		slog.Warn("Create new elasticsearch client failed", err)
		return nil, err
	}
	_, err = esClient.Info()
	if err != nil {
		slog.Warn("Create new elasticsearch client INFO failed", err)
		return nil, err
	}
	return &elasticSearch{
		Client: esClient,
	}, nil
}

// Search implements ElasticSearch.
func (es *elasticSearch) Search(indexName, body string) (map[string]interface{}, error) {
	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(body),
	}
	// Search
	res, err := req.Do(context.Background(), es.Client)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		errTxt := fmt.Sprintf("Elasticsearch search error: %s", res.Status())
		return nil, errors.New(errTxt)
	}

	result := make(map[string]interface{})

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		slog.Error("decode-body failed", err)
		return nil, err
	}
	return result, nil
}

// Insert implements ElasticSearch.
func (es *elasticSearch) Insert(index string, data any, documentID string) error {
	out, err := json.Marshal(data)
	if err != nil {
		slog.Warn("Error marshalling", err)
		return err
	}
	req := esapi.IndexRequest{
		DocumentID: documentID,
		Index:      index,
		Body:       strings.NewReader((string(out))),
		Refresh:    "true",
	}
	//Insert into elastic
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		slog.Warn("Error inserting elastic", err)
		return err
	}
	defer res.Body.Close()
	response := ResponseRequest{}

	if res.IsError() {
		var e ResponseError
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return err
		} else {
			if e.Error.Reason != "" {
				errCus := errors.New(e.Error.Reason + "->" + e.Error.CausedBy.Reason)
				slog.Warn("Error inserting elastic", errCus)
				return errCus
			} else {
				return errors.New("response error elasticsearch invalid, can't find reason")
			}
		}
	} else {
		err := json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		} else {
			if strings.ToLower(response.Result) == "created" {
				return nil
			}
			return errors.New("not inserted")
		}
	}
}

// Ping implements ElasticSearch.
func (es *elasticSearch) Ping() error {
	// Perform a ping to check the Elasticsearch cluster's availability
	res, err := es.Client.Ping()
	if err != nil {
		slog.Warn("Ping Elasticsearch failed", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		slog.Warn("Error to connect with Elasticsearch", err)
		return errors.New("error to connect with elasticsearch")
	} else {
		return nil
	}
}

// Remove implements ElasticSearch.
func (es *elasticSearch) Remove(documentID string, index string) error {
	req := esapi.DeleteRequest{
		DocumentID: documentID,
		Index:      index,
		Refresh:    "true",
	}
	// Do Request Remove Item
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		slog.Warn("Remove Item Elasticsearch failed", err)
		return err
	}
	defer res.Body.Close()
	return nil
}
