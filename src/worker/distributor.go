package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail,	opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}



func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail,	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to encode task payload %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts... )
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("Failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()). 
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil 
}


func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
		var payload PayloadSendVerifyEmail 
		if err := json.Unmarshal(task.Payload(), &payload); err != nil  {
			return fmt.Errorf("Failed to unmarshal payload: %w", asynq.SkipRetry)
		}

		user, err := processor.store.GetUser(ctx, payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				// return fmt.Errorf("user doesn't exist: %v", asynq.SkipRetry)
				return fmt.Errorf("user doesn't exist: %v", err)
			}
			return fmt.Errorf("Failed to get user: %w", err)
		}

		verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
			Username: payload.Username,
			Email: user.Email,
			SecretCode: utils.RandomString(32),
		})
		if err != nil {
			return fmt.Errorf("Failed to create verify email: %v", err)
		}


		verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)

		subject := "Welcome to shiny bank"
		content := fmt.Sprintf(`Hello %s, <br/>
		Thank you for registering with us! <br/>
		Please <a href="%s"> click here</a> to verify your email address
		`, user.FullName, verifyUrl)

		to := []string{user.Email}

		err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Failed to send verify email: %w", err)
		}

		// TODO: send email to user 
		log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()). 
			Str("email", user.Email).Msg("Processor task")

		return nil 
}