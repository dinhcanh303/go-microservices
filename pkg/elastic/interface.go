package elastic

type ElasticSearch interface {
	Insert(index string, data any, documentID string) error
	Remove(documentID string, index string) error
	Ping() error
	Search(indexName, body string) (map[string]interface{}, error)
}
