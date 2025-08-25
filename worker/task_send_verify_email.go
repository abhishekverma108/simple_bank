package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"go.elastic.co/apm/v2"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	span, _ := apm.StartSpan(ctx, "asynq.enqueue_send_verify_email", "taskqueue")
	defer span.End()

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return fmt.Errorf("failed to marshal task payload %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return fmt.Errorf("failed to enqueue task %w", err)

	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	span, _ := apm.StartSpan(ctx, "process_send_verify_email", "taskqueue")
	defer span.End()

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		errWrapped := fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
		apm.CaptureError(ctx, errWrapped).Send()
		return errWrapped
	}

	spanDB, _ := apm.StartSpan(ctx, "db.get_user_by_username", "db")
	user, err := processor.store.GetUserByUsername(ctx, payload.Username)
	spanDB.End()
	if err != nil {
		if err == sql.ErrNoRows {
			errWrapped := fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
			apm.CaptureError(ctx, errWrapped).Send()
			return errWrapped
		}
		errWrapped := fmt.Errorf("failed to get user: %w", err)
		apm.CaptureError(ctx, errWrapped).Send()
		return errWrapped
	}
	// --- write username/email marker to Redis ---
	key := fmt.Sprintf("verify_email:%s", payload.Username)
	spanRedis, _ := apm.StartSpan(ctx, "redis.set_verify_email_marker", "cache")
	if err := processor.redisClient.Set(ctx, key, payload.Username, 0).Err(); err != nil {
		spanRedis.End()
		errWrapped := fmt.Errorf("failed to set redis key %q: %w", key, err)
		apm.CaptureError(ctx, errWrapped).Send()
		return errWrapped
	}
	spanRedis.End()
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")
	return nil
}
