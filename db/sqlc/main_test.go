package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
	apmsql "go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
)

const (
	dbDriver = "postgres"
	// dbSource = "postgresql://postgres:pae9bai7Cahg?ahcae"g@134.209.150.195:5445/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = apmsql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
