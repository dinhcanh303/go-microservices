package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/go-ldap/ldap/v3"
)

type ldapClient struct {
	attributes         []string
	config             *configs.Ldap
	conn               *ldap.Conn
	clientCertificates []tls.Certificate
}

// Authenticate implements LdapClient.
func (l *ldapClient) Authenticate(username string, password string) (bool, map[string]string, error) {
	err := l.Connect()
	if err != nil {
		return false, nil, err
	}
	if l.config.LdapBindDN != "" && l.config.LdapPassword != "" {
		err := l.conn.Bind(l.config.LdapBindDN, l.config.LdapPassword)
		if err != nil {
			return false, nil, err
		}
	}
	attributes := append(l.attributes, "dn")
	filter := fmt.Sprintf("(%s=%s)", l.config.LdapFilter, username)
	searchRequest := ldap.NewSearchRequest(
		l.config.LdapBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		l.config.LdapTimeout,
		false,
		filter,
		attributes,
		nil,
	)
	sr, err := l.conn.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}
	if len(sr.Entries) < 1 {
		return false, nil, errors.New("user does not exist")
	}
	if len(sr.Entries) > 1 {
		return false, nil, errors.New("too many entries returned")
	}
	userDN := sr.Entries[0].DN
	err = l.conn.Bind(userDN, password)
	if err != nil {
		return false, nil, err
	}
	user := map[string]string{}
	for _, attr := range l.attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}
	return true, user, nil
}

// Close implements LdapClient.
func (l *ldapClient) Close() {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}

// Connect implements LdapClient.
func (l *ldapClient) Connect() error {
	if l.conn == nil {
		var ldapConn *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", l.config.LdapHost, l.config.LdapPort)
		if !l.config.LdapSSL {
			ldapConn, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}
			if !l.config.LdapTLS {
				err = l.conn.StartTLS(&tls.Config{
					InsecureSkipVerify: true,
				})
				if err != nil {
					return err
				}
			}
		} else {
			config := &tls.Config{
				InsecureSkipVerify: l.config.LdapTLS,
				ServerName:         l.config.LdapBindDN,
			}
			if l.clientCertificates != nil && len(l.clientCertificates) > 0 {
				config.Certificates = l.clientCertificates
			}
			ldapConn, err = ldap.DialTLS("tcp", address, config)
			if err != nil {
				return err
			}
		}
		l.conn = ldapConn
	}
	return nil
}

var _ LdapClient = (*ldapClient)(nil)
