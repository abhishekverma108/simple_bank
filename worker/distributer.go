package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}
type RedisTaskDistributor struct {
	client      *asynq.Client
	redisClient *redis.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt, redisClient *redis.Client) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client:      client,
		redisClient: redisClient,
	}

}

// GetRedisClient returns the monitored Redis client for direct operations
func (distributor *RedisTaskDistributor) GetRedisClient() *redis.Client {
	return distributor.redisClient
}
