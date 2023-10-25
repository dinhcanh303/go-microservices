package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Minio struct {
	EndPoint        string `env-required:"true" yaml:"end_point" env:"MINIO_ENDPOINT"`
	AccessKeyID     string `env-required:"true" yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `env-required:"true" yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
	Location        string `env-required:"true" yaml:"location" env:"MINIO_DEFAULT_REGION"`
	BucketName      string `env-required:"true" yaml:"bucket_name" env:"MINIO_BUCKET"`
	RootFolder      string `env-required:"true" yaml:"root_folder" env:"MINIO_ROOT_FOLDER"`
	UseSSL          bool   `env-required:"true" yaml:"use_ssl" env:"MINIO_USE_SSL"`
}

func NewConfigMinio() (*Minio, error) {
	cfg := &Minio{}

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
