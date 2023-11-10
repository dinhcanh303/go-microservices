package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type elasticSearch struct {
	Client *elasticsearch.Client
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

var ES ElasticSearch = (*elasticSearch)(nil)

func NewElasticSearch(cfg configs.ElasticSearch) (ElasticSearch, error) {
	address := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	cfgES := elasticsearch.Config{
		Addresses: []string{address},
		Username:  cfg.UserName,
		Password:  cfg.Password,
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
