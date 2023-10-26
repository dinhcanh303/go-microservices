package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Minio struct {
	EndPoint        string `env:"MINIO_ENDPOINT" env-default:"minio_endpoint"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID" env-default:"access_key_id"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" env-default:"secret_access_key"`
	Region          string `env:"MINIO_DEFAULT_REGION" env-default:"region"`
	BucketName      string `env:"MINIO_BUCKET" env-default:"bucket"`
	RootFolder      string `env:"MINIO_ROOT_FOLDER" env-default:"folder"`
	UseSSL          bool   `env:"MINIO_USE_SSL" env-default:"false"`
}

func NewConfigMinio() (*Minio, error) {
	cfg := &Minio{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
