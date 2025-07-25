package handler

import (
	"errors"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/dto"
	domainErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/response"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
	cfg *config.Config
}

func NewAuthHandler(svc *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg}
}

func (h *AuthHandler) JWKS(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "application/json", h.svc.JWKS())
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err)
		return
	}

	tokens, err := h.svc.RegisterUser(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserAlreadyExists):
			response.RespondWithError(ctx, http.StatusConflict,
				"пользователь с таким email уже зарегистрирован", nil)
			return
		default:
			response.RespondWithError(ctx, http.StatusInternalServerError,
				"внутренняя ошибка сервера", nil)
			return
		}
	}

	ctx.JSON(http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err)
		return
	}
	tokens, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidCredentials):
			response.RespondWithError(ctx, http.StatusUnauthorized,
				"неверный email или пароль", nil)
			return
		default:
			response.RespondWithError(ctx, http.StatusInternalServerError,
				"внутренняя ошибка сервера", nil)
			return
		}
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err)
		return
	}
	tokens, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidToken):
			response.RespondWithError(ctx, http.StatusUnauthorized,
				"неверный токен", nil)
			return
		default:
			response.RespondWithError(ctx, http.StatusInternalServerError,
				"внутренняя ошибка сервера", nil)
			return
		}
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) ClientToken(ctx *gin.Context) {
	var req dto.ClientCredentialsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleValidationError(ctx, err)
		return
	}
	tokens, err := h.svc.ClientToken(ctx.Request.Context(), req.ClientID, req.ClientSecret)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidCredentials):
			response.RespondWithError(ctx, http.StatusUnauthorized,
				"неверные учетные данные клиента", nil)
			return
		default:
			response.RespondWithError(ctx, http.StatusInternalServerError,
				"внутренняя ошибка сервера", nil)
			return
		}
	}
	ctx.JSON(http.StatusOK, tokens)
}
