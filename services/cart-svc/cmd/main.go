package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/cart-svc/config"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/authclient"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/controller"
	middleware "github.com/oidiral/e-commerce/services/cart-svc/internal/controller/middleware"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/db"
	grpcx "github.com/oidiral/e-commerce/services/cart-svc/internal/grpc"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/pb/catalog"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/repository/postgres"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
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
	logger.Info().Msg("Cart service started")
	database, err := db.NewPool(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()
	logger.Info().Msg("Connected to database")

	redisClient, err := db.NewRedisClient(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to redis")
	}

	authClie := authclient.NewClient(cfg.AuthURL, cfg.ClientID, cfg.ClientSecret)
	cartRepo := postgres.NewCartRepoPg(database)
	conn, err := grpc.NewClient(cfg.CatalogGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(grpcx.NewJWTPerRPCCreds(authClie)))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to catalog service")
	}
	catalogClient := catalog.NewCatalogClient(conn)
	cartService := service.NewCartService(cartRepo, logger, redisClient, catalogClient)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))

	controller.RegisterRoutes(router, cartService)

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
