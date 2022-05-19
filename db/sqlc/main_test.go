package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	dbSource = "jija:password@/simple_bank?parseTime=true"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("database created, but cannot be pinged", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
