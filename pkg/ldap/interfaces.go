package ldap

type LdapClient interface {
	Close()
	Authenticate(username, password string) (bool, map[string]string, error)
	Connect() error
}
