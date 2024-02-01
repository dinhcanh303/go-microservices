package meili

import (
	"github.com/google/wire"
	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearchConn struct {
	Host   string
	APIKey string
}
type meiliSearch struct {
	Client *meilisearch.Client
}

var _ MeiliSearch = (*meiliSearch)(nil)

var MeiliSearchSet = wire.NewSet(NewMeiliSearch)

func NewMeiliSearch(meiliSearchConn MeiliSearchConn) MeiliSearch {
	cfg := meilisearch.ClientConfig{
		Host:   meiliSearchConn.Host,
		APIKey: meiliSearchConn.APIKey,
	}
	client := meilisearch.NewClient(cfg)
	return &meiliSearch{
		Client: client,
	}
}
func (m *meiliSearch) UpdateTypoTolerance(index string, config meilisearch.MinWordSizeForTypos) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).UpdateTypoTolerance(
		&meilisearch.TypoTolerance{
			MinWordSizeForTypos: config,
		},
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *meiliSearch) UpdateRankingRules(index string, rankingRules *[]string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).UpdateRankingRules(rankingRules)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (m *meiliSearch) DeleteIndex(index string, primaryKey ...string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.DeleteIndex(index)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *meiliSearch) AddDocuments(index string, documents interface{}, primaryKey ...string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).AddDocuments(documents, primaryKey...)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *meiliSearch) UpdateFilterableAttributes(index string, req *[]string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).UpdateFilterableAttributes(req)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *meiliSearch) UpdateSearchableAttributes(index string, req *[]string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).UpdateSearchableAttributes(req)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *meiliSearch) Search(index, search string, configs *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	if configs == nil {
		configs = &meilisearch.SearchRequest{
			Limit:  10,
			Offset: 0,
		}
	}
	searchRes, err := m.Client.Index(index).Search(search, configs)
	if err != nil {
		return nil, err
	}
	return searchRes, err
}
func (m *meiliSearch) DeleteDocuments(index string, ids []string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).DeleteDocuments(ids)
	if err != nil {
		return nil, err
	}
	return task, err
}
func (m *meiliSearch) DeleteDocument(index string, id string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).DeleteDocument(id)
	if err != nil {
		return nil, err
	}
	return task, err
}
func (m *meiliSearch) DeleteAllDocuments(index string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.Index(index).DeleteAllDocuments()
	if err != nil {
		return nil, err
	}
	return task, err
}
func (m *meiliSearch) CreateIndex(index string, primaryKey string) (*meilisearch.TaskInfo, error) {
	task, err := m.Client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        index,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		return nil, err
	}
	return task, err
}
