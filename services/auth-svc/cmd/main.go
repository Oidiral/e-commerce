package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/db"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/middleware"
	"os"
	"time"

	"github.com/oidiral/e-commerce/services/auth-svc/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := SetupLogger(cfg.Env)

	logger.Info().Msg("Auth service started")

	database, err := db.NewPool(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()
	logger.Info().Msg("Connected to database")

	r := gin.New()
	r.Use(gin.Recovery(), middleware.ApiErrorMiddleware(), gin.Logger())

}

func SetupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	fmt.Println("Logger environment:", env)

	switch env {
	case envLocal:
		return log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}).Level(zerolog.DebugLevel)
	case envDev:
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Level(zerolog.DebugLevel)
	case envProd:
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger().
			Level(zerolog.InfoLevel)
	default:
		panic("unknown environment: " + env)
	}
}
