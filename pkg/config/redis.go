package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Redis struct {
	Host     string `env:"REDIS_HOST" env-default:"redis"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"password"`
}

func NewConfigRedis() (*Redis, error) {
	cfg := &Redis{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
