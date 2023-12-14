package config

import (
	"fmt"
	"log"
	"os"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		configs.App   `yaml:"app"`
		configs.HTTP  `yaml:"http"`
		configs.Log   `yaml:"logger"`
		RabbitMQ      `yaml:"rabbitmq"`
		ElasticSearch `yaml:"elasticsearch"`
		MeiliSearch   `yaml:"meilisearch"`
	}
	RabbitMQ struct {
		URL string `env-required:"true" yaml:"url" env:"RABBITMQ_URL"`
	}
	ElasticSearch struct {
		URL      string `env-required:"true" yaml:"url" env:"ES_URL"`
		UserName string `env-required:"true" yaml:"username" env:"ES_USERNAME"`
		Password string `env-required:"true" yaml:"password" env:"ES_PASSWORD"`
	}
	MeiliSearch struct {
		Host   string `env-required:"true" yaml:"host" env:"ML_HOST"`
		ApiKey string `env-required:"true" yaml:"api_key" env:"ML_API_KEY"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println("config path: " + dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
