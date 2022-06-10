package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cyberdr0id/bank-backend/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	if err := testDB.Ping(); err != nil {
		log.Fatal("database created, but cannot be pinged", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
