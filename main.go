package main

import (
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := apmsql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
