package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	user "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	AppErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	repository "github.com/oidiral/e-commerce/services/auth-svc/internal/repository"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/utils"
	"github.com/rs/zerolog"
	"time"
)

type AuthService struct {
	repo repository.AuthRepository
	log  zerolog.Logger
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, log zerolog.Logger, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, log: log, cfg: cfg}
}

type TokenPair struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	AccessExpiresAt  int64  `json:"accessExpiresAt"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string) (*TokenPair, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		s.log.Error().Err(err).Msg("hash pwd")
		return nil, err
	}

	u, err := s.repo.CreateIfNotExists(ctx, email, hash)
	if err != nil {
		return nil, err
	}

	access, accExp, err := s.createAccessToken(u)
	if err != nil {
		return nil, err
	}
	refresh, refExp, err := s.createRefreshToken(u)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresAt:  accExp.Unix(),
		RefreshExpiresAt: refExp.Unix(),
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if !token.Valid {
		s.log.Error().Msg("invalid refresh token")
		return nil, AppErr.ErrInvalidToken
	}
	if err != nil {
		s.log.Error().Err(err).Msg("parse refresh token")
		return nil, AppErr.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.log.Error().Msg("invalid token claims")
		return nil, AppErr.ErrInvalidToken
	}
	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		s.log.Error().Msg("refresh token expired")
		return nil, AppErr.ErrInvalidToken
	}
	subStr, ok := claims["sub"].(string)
	subUUID, err := uuid.Parse(subStr)
	if err != nil {
		s.log.Error().Err(err).Msg("invalid token subject format")
		return nil, AppErr.ErrInvalidToken
	}
	if !ok {
		s.log.Error().Msg("invalid token subject")
		return nil, AppErr.ErrInvalidToken
	}
	email, ok := claims["email"].(string)
	if !ok {
		s.log.Error().Msg("invalid token email")
		return nil, AppErr.ErrInvalidToken
	}
	role, ok := claims["role"].(string)
	if !ok {
		s.log.Error().Msg("invalid token role")
		return nil, AppErr.ErrInvalidToken
	}

	u := user.User{
		ID:    subUUID,
		Email: email,
		Role:  role,
	}
	access, accExp, err := s.createAccessToken(&u)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create access token")
		return nil, err
	}

	return &TokenPair{
		AccessToken:      access,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accExp.Unix(),
		RefreshExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0).Unix(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, AppErr.ErrNotFound) {
			return nil, AppErr.ErrInvalidCredentials
		}
		s.log.Error().Err(err).Msg("get user by email")
		return nil, err
	}

	if !utils.CheckPasswordHash(password, u.Password) {
		return nil, AppErr.ErrInvalidCredentials
	}

	access, accExp, err := s.createAccessToken(u)
	if err != nil {
		return nil, err
	}
	refresh, refExp, err := s.createRefreshToken(u)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresAt:  accExp.Unix(),
		RefreshExpiresAt: refExp.Unix(),
	}, nil
}

func (s *AuthService) createAccessToken(u *user.User) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.cfg.JWT.AccessTokenTTL)
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  u.Role,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
	jw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		s.log.Error().Err(err).Msg("create access token")
		return "", time.Time{}, err
	}
	return jw, exp, nil
}

func (s *AuthService) createRefreshToken(u *user.User) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.cfg.JWT.RefreshTokenTTL)
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  u.Role,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
	jw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		s.log.Error().Err(err).Msg("create refresh token")
		return "", time.Time{}, err
	}
	return jw, exp, nil
}
