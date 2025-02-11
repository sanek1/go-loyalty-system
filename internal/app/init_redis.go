package app

import (
	"context"
	"go-loyalty-system/pkg/logging"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func initRedis(ctx context.Context, l *logging.ZapLogger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		l.FatalCtx(ctx, "app - initRedis - client.Ping: %w", zap.Error(err))
	}

	return client
}
