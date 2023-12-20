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
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
		Mongo        `yaml:"mongodb"`
	}

	Mongo struct {
		MongoURL string `env-required:"true" yaml:"mongo_url" env:"MONGO_URL"`
		DBName   string `env-required:"true" yaml:"db_name" env:"MONGO_DB_NAME"`
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
