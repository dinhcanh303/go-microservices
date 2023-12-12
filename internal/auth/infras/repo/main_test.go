package repo

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/dinhcanh303/go-microservices/internal/auth/infras/postgresql"
	_ "github.com/lib/pq"
)

const (
	dbDriver  = "postgres"
	dbSource  = "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable"
	dbSource2 = "postgres://postgres:123456@127.0.0.1:5433/postgres?sslmode=disable"
)

var testQueries *postgresql.Queries
var testRepQueries *postgresql.Queries

func TestMain(m *testing.M) {
	dbM, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}
	dbR, err := sql.Open(dbDriver, dbSource2)
	if err != nil {
		log.Fatal("can't connect to db rep:", err)
	}
	testQueries = postgresql.New(dbM)
	testRepQueries = postgresql.New(dbR)
	os.Exit(m.Run())
}
