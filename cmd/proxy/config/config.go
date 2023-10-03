package config

import (
	"fmt"
	configs "go-microservices/pkg/config"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
		GRPC         `yaml:"grpc"`
	}
	GRPC struct {
		GroupHost string `env-required:"true" yaml:"group_host" env:"GRPC_GROUP_HOST"`
		GroupPort int    `env-required:"true" yaml:"group_port" env:"GRPC_GROUP_PORT"`
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
