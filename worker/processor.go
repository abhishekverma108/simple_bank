package worker

import (
	"context"
	db "simplebank/db/sqlc"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server      *asynq.Server
	store       db.Store
	redisClient *redis.Client
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, redisClient *redis.Client, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{},
	)
	return &RedisTaskProcessor{
		server:      server,
		redisClient: redisClient,
		store:       store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
