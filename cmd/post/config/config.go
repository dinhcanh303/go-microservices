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
		PG            `yaml:"postgres"`
		CommentClient `yaml:"comment_client"`
		LikeClient    `yaml:"like_client"`
		UploadClient  `yaml:"upload_client"`
		GroupClient   `yaml:"group_client"`
		AuthClient    `yaml:"auth_client"`
	}

	PG struct {
		PoolMax  int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DbURL    string `env-required:"true" yaml:"db_url" env:"PG_URL"`
		DbRepURL string `env-required:"true" yaml:"db_rep_url" env:"PG_REP_URL"`
	}

	CommentClient struct {
		URL string `env-required:"true" yaml:"comment_url" env:"COMMENT_CLIENT_URL"`
	}
	LikeClient struct {
		URL string `env-required:"true" yaml:"like_url" env:"LIKE_CLIENT_URL"`
	}
	UploadClient struct {
		URL string `env-required:"true" yaml:"upload_url" env:"UPLOAD_CLIENT_URL"`
	}
	GroupClient struct {
		URL string `env-required:"true" yaml:"group_url" env:"GROUP_CLIENT_URL"`
	}
	AuthClient struct {
		URL string `env-required:"true" yaml:"auth_url" env:"AUTH_CLIENT_URL"`
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
