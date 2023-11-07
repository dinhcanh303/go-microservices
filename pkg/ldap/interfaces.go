package ldap

type LdapClient interface {
	Connect() error
	Close()
	Authenticate(username, password string) (bool, map[string]string, error)
}
