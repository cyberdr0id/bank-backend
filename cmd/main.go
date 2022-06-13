package main

import (
	"database/sql"
	"log"

	"github.com/cyberdr0id/bank-backend/api"
	db "github.com/cyberdr0id/bank-backend/db/sqlc"
	"github.com/cyberdr0id/bank-backend/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unbale to load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server", err)
	}
}
