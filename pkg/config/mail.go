package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Mail struct {
	Host        string `env:"MAIL_HOST" env-default:"smtp.email.com"`
	Port        string `env:"MAIL_PORT" env-default:"587"`
	Password    string `env:"MAIL_PASSWORD" env-default:"password"`
	Encryption  string `env:"MAIL_ENCRYPTION" env-default:"tls"`
	FromAddress string `env:"MAIL_FROM_ADDRESS" env-default:"example@example.com"`
	FromName    string `env:"MAIL_FROM_NAME" env-default:"example"`
	Bcc         string `env:"MAIL_BCC" env-default:"dinhcanhng303@gmail.com"`
}

func NewConfigMail() (*Mail, error) {
	cfg := &Mail{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
