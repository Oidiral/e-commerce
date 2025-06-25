package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/middleware"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/service"
)

func RegisterRoutes(router *gin.Engine, authService *service.AuthService, cfg *config.Config) {
	router.Use(middleware.ApiErrorMiddleware())

	authHandler := NewAuthHandler(authService, cfg)
	router.GET("/.well-known/jwks.json", authHandler.JWKS)
	api := router.Group("/api/v1/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)
		api.POST("/token", authHandler.ClientToken)
	}
}
