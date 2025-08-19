package main

import (
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"simplebank/worker"

	"github.com/hibiken/asynq"
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
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	go runTaskProcessor(redisOpts, store)
	server := api.NewServer(store, taskDistributor)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Println("starting task processor...")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("cannot start task processor:", err)
	}
}
