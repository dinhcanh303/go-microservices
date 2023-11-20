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
		PG           `yaml:"postgres"`
		GRPC         `yaml:"grpc"`
	}
	GRPC struct {
		GroupHost   string `env-required:"true" yaml:"group_host" env:"GRPC_GROUP_HOST"`
		GroupPort   int    `env-required:"true" yaml:"group_port" env:"GRPC_GROUP_PORT"`
		PostHost    string `env-required:"true" yaml:"post_host" env:"GRPC_POST_HOST"`
		PostPort    int    `env-required:"true" yaml:"post_port" env:"GRPC_POST_PORT"`
		CommentHost string `env-required:"true" yaml:"comment_host" env:"GRPC_COMMENT_HOST"`
		CommentPort int    `env-required:"true" yaml:"comment_port" env:"GRPC_COMMENT_PORT"`
		LikeHost    string `env-required:"true" yaml:"like_host" env:"GRPC_LIKE_HOST"`
		LikePort    int    `env-required:"true" yaml:"like_port" env:"GRPC_LIKE_PORT"`
		UploadHost  string `env-required:"true" yaml:"upload_host" env:"GRPC_UPLOAD_HOST"`
		UploadPort  int    `env-required:"true" yaml:"upload_port" env:"GRPC_UPLOAD_PORT"`
		AuthHost    string `env-required:"true" yaml:"auth_host" env:"GRPC_AUTH_HOST"`
		AuthPort    int    `env-required:"true" yaml:"auth_port" env:"GRPC_AUTH_PORT"`
	}
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DsnURL  string `env-required:"true" yaml:"dsn_url" env:"PG_DSN_URL"`
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
