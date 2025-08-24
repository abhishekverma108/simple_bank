package main

import (
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"simplebank/worker"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8/v2"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
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
	// Create monitored Redis client for both asynq and direct operations
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
	})
	redisClient.AddHook(apmgoredis.NewHook())
	// Configure Redis options for asynq
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts, redisClient)
	go runTaskProcessor(redisOpts, redisClient, store)
	server := api.NewServer(store, taskDistributor, redisClient)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
func runTaskProcessor(redisOpt asynq.RedisClientOpt, redisClient *redis.Client, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, redisClient, store)
	log.Println("starting task processor...")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("cannot start task processor:", err)
	}
}
