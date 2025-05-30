package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
)

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode,
	)

	c, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse c: %w", err)
	}

	c.MaxConns = int32(cfg.Database.MaxOpenConns)
	c.MinConns = int32(cfg.Database.MaxIdleConns)
	c.ConnConfig.ConnectTimeout = cfg.Database.ConnTimeout

	pool, err := pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.ConnTimeout)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
