package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Ldap struct {
	LdapLogging      bool   `env:"LDAP_LOGGING" env-default:"true"`
	LdapConnection   string `env:"LDAP_CONNECTION" env-default:"default"`
	LdapHost         string `env:"LDAP_HOST" env-default:"localhost"`
	LdapBindDN       string `env:"LDAP_BIND_DN" env-default:"dc=example,dc=com"`
	LdapPassword     string `env:"LDAP_PASSWORD" env-default:"password"`
	LdapPort         int    `env:"LDAP_PORT" env-default:"389"`
	LdapBaseDN       string `env:"LDAP_BASE_DN" env-default:"dc=example,dc=com"`
	LdapTimeout      int    `env:"LDAP_TIMEOUT" env-default:"10"`
	LdapSSL          bool   `env:"LDAP_SSL" env-default:"false"`
	LdapTLS          bool   `env:"LDAP_TLS" env-default:"false"`
	LdapFilter       string `env:"LDAP_FILTER" env-default:"username"`
	LdapUserNameTest string `env:"LDAP_USERNAME_TEST" env-default:"username_test"`
	LdapPasswordTest string `env:"LDAP_PASSWORD_TEST" env-default:"password_test"`
}

func NewLdapConfig() (*Ldap, error) {
	cfg := &Ldap{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
