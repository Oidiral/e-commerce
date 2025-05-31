package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/dto"
	domainErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/response"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
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
