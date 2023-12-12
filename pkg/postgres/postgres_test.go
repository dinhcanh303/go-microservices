package postgres

import (
	"testing"

	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

const (
	DnsURL = "host=127.0.0.1 user=postgres password=123456 dbname=postgres sslmode=disable"
)

func TestPostgres(t *testing.T) {
	err := utils.LoadFileEnvOnLocal()
	require.NoError(t, err)
	dbEngine, err := NewPostgresDB(DBConnString(DnsURL))
	require.NoError(t, err)
	require.NotEmpty(t, dbEngine)
}
