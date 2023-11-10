package configs

import "github.com/ilyakaznacheev/cleanenv"

type ElasticSearch struct {
	Host     string `env:"ES_HOST" env-default:"localhost"`
	Port     string `env:"ES_PORT" env-default:"9200"`
	UserName string `env:"ES_USERNAME" env-default:"elastic"`
	Password string `env:"ES_PASSWORD" env-default:"elastic"`
}

func NewConfigElasticSearch() (*ElasticSearch, error) {
	cfg := &ElasticSearch{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
