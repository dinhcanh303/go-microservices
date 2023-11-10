package ldap

import (
	"testing"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestLdap(t *testing.T) {
	err := utils.LoadFileEnvOnLocal()
	require.NoError(t, err)
	configLdap, err := configs.NewLdapConfig()
	require.NoError(t, err)
	require.NotEmpty(t, configLdap)
	ldapClient := NewLdapClient(configLdap, []string{configLdap.LdapFilter})
	require.NotEmpty(t, ldapClient)
	username := utils.RandomEmailCompany()
	password := utils.RandomString(8)
	check, _, err := ldapClient.Authenticate(username, password)
	require.EqualError(t, err, "user does not exist")
	require.Equal(t, check, false)
	username = configLdap.LdapUserNameTest
	password = configLdap.LdapPasswordTest
	check, _, err = ldapClient.Authenticate(username, password)
	require.NoError(t, err)
	require.Equal(t, check, true)
}
