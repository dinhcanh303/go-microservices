package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// DistributeTaskSendVerifyEmail implements TaskDistributor.
func (*RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, opts ...asynq.Option) error {
	panic("unimplemented")
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
