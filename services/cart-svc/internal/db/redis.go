package db

import (
	"context"
	"fmt"
	"github.com/oidiral/e-commerce/services/cart-svc/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	ctx := context.Background()

	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rds.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisClient{Client: rds}, nil
}
