package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://admin:pass12345@localhost:5432/simple_bank?sslmode=disable"
)

var testQuerys *Queries

func TestMain(m *testing.M) {
	cnn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connet to db:", err)
	}

	testQuerys = New(cnn)

	os.Exit(m.Run())
}
