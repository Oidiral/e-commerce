package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/db"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/handler"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/middleware"
	repository "github.com/oidiral/e-commerce/services/auth-svc/internal/repository"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/service"
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

	authRepo := repository.NewAuthRepository(database)
	logger.Info().Msg("Auth repository initialized")
	authService := service.NewAuthService(authRepo, logger, cfg)
	logger.Info().Msg("Auth service initialized")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	handler.RegisterRoutes(router, authService)
	logger.Info().Msg("Routes registered")
	if err := router.Run(cfg.Server.Port); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func SetupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	fmt.Println("Logger environment:", env)

	switch env {
	case envLocal:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	case envDev:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	case envProd:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	default:
		panic("unknown environment: " + env)
	}
}
