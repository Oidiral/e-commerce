package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/db"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/handler"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/middleware"
	repository "github.com/oidiral/e-commerce/services/auth-svc/internal/repository"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ClientRepo := repository.NewClientRepository(database)
	authRepo := repository.NewAuthRepository(database)
	logger.Info().Msg("Auth repository initialized")
	authService := service.NewAuthService(authRepo, logger, cfg, ClientRepo)
	logger.Info().Msg("Auth service initialized")
	logger.Info().Msg("db_name" + ": " + cfg.Database.Name)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	handler.RegisterRoutes(router, authService, cfg)
	logger.Info().Msg("Routes registered")
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown failed")
	} else {
		logger.Info().Msg("Server gracefully stopped")
	}
}

func SetupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Info().
		Str("env", env).
		Msg("Logger environment")

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
