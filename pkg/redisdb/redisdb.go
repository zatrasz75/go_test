package redisdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"zatrasz75/go_test/pkg/logger"
)

type Redis struct {
	connAttempts int
	connTimeout  time.Duration

	Rds        *redis.Client
	Expiration time.Duration
}

func New(addr string, l logger.LoggersInterface, opts ...Option) (*Redis, error) {
	r := &Redis{}

	r.Rds = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	// Пользовательские параметры
	for _, opt := range opts {
		opt(r)
	}

	// Проверка доступности Redis
	for i := 0; i < r.connAttempts; i++ {
		_, err := r.Rds.Ping(context.Background()).Result()
		if err == nil {
			break
		}
		l.Info("Redis пытается подключиться, попыток осталось: %d", r.connAttempts-i)

		time.Sleep(r.connTimeout)
	}

	if _, err := r.Rds.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - не удалось подключиться: %w", err)
	}

	return r, nil
}
