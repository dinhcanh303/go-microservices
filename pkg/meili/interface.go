package meili

import "github.com/meilisearch/meilisearch-go"

type MeiliSearch interface {
	AddDocuments(index string, documents interface{}, primaryKey ...string) (*meilisearch.TaskInfo, error)
	DeleteDocuments(index string, ids []string) (*meilisearch.TaskInfo, error)
	DeleteDocument(index string, id string) (*meilisearch.TaskInfo, error)
	DeleteAllDocuments(index string) (*meilisearch.TaskInfo, error)
	UpdateFilterableAttributes(index string, req *[]string) (*meilisearch.TaskInfo, error)
	UpdateSearchableAttributes(index string, req *[]string) (*meilisearch.TaskInfo, error)
	Search(index, search string, configs *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error)
	CreateIndex(index string, primaryKey string) (*meilisearch.TaskInfo, error)
}
