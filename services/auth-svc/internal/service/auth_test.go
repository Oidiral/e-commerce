package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	model "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	AppErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/repository/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()

	password := "testpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	userId := uuid.New()
	mockUser := &model.User{
		ID:       userId,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Roles:    []string{"user"},
		Status:   1,
	}

	mockRepo.
		EXPECT().
		GetByEmail(gomock.Any(), gomock.Eq("test@example.com")).
		Return(mockUser, nil)

	authService := NewAuthService(mockRepo, logger, cfg)

	resp, loginErr := authService.Login(context.Background(), "test@example.com", password)

	assert.NoError(t, loginErr)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken, "должен быть сгенерирован AccessToken")
	assert.NotEmpty(t, resp.RefreshToken, "должен быть сгенерирован RefreshToken")
	assert.True(t, resp.AccessExpiresAt > time.Now().Unix(), "AccessExpiresAt > now")
	assert.True(t, resp.RefreshExpiresAt > resp.AccessExpiresAt, "RefreshExpiresAt > AccessExpiresAt")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()

	realHash, _ := bcrypt.GenerateFromPassword([]byte("correctPassword"), bcrypt.DefaultCost)
	mockUser := &model.User{
		ID:       uuid.New(),
		Email:    "test2@example.com",
		Password: string(realHash),
		Roles:    []string{"user"},
		Status:   1,
	}

	mockRepo.
		EXPECT().
		GetByEmail(gomock.Any(), gomock.Eq("test2@example.com")).
		Return(mockUser, nil)

	authService := NewAuthService(mockRepo, logger, cfg)

	_, err := authService.Login(context.Background(), "test2@example.com", "wrongPassword")

	assert.ErrorIs(t, err, AppErr.ErrInvalidCredentials)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()

	mockRepo.
		EXPECT().
		GetByEmail(gomock.Any(), gomock.Eq("nouser@example.com")).
		Return(nil, AppErr.ErrNotFound)

	authService := NewAuthService(mockRepo, logger, cfg)

	_, err := authService.Login(context.Background(), "nouser@example.com", "anyPassword")

	assert.ErrorIs(t, err, AppErr.ErrInvalidCredentials)
}

func TestAuthService_RegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()

	userId := uuid.New()
	mockUser := &model.User{
		ID:    userId,
		Email: "newuser@example.com",
		Roles: []string{"user"},
	}

	mockRepo.
		EXPECT().
		CreateIfNotExists(gomock.Any(), "newuser@example.com", gomock.Any()).
		Return(mockUser, nil)

	authService := NewAuthService(mockRepo, logger, cfg)

	resp, err := authService.RegisterUser(context.Background(), "newuser@example.com", "plainPassword")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.True(t, resp.AccessExpiresAt > time.Now().Unix())
	assert.True(t, resp.RefreshExpiresAt > resp.AccessExpiresAt)
}

func TestAuthService_RegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()

	mockRepo.
		EXPECT().
		CreateIfNotExists(gomock.Any(), "existing@example.com", gomock.Any()).
		Return(nil, AppErr.ErrUserAlreadyExists)

	authService := NewAuthService(mockRepo, logger, cfg)

	_, err := authService.RegisterUser(context.Background(), "existing@example.com", "anyPassword")

	assert.ErrorIs(t, err, AppErr.ErrUserAlreadyExists)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()
	authService := NewAuthService(nil, logger, cfg)

	userId := uuid.New()
	claims := jwt.MapClaims{
		"sub":   userId.String(),
		"email": "user@example.com",
		"roles": []string{"user"},
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := tokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	resp, err := authService.Refresh(refreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, refreshToken, resp.RefreshToken)
	assert.True(t, resp.AccessExpiresAt > time.Now().Unix())
}

func TestAuthService_Refresh_ExpiredToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	logger := zerolog.Nop()
	authService := NewAuthService(nil, logger, cfg)

	claims := jwt.MapClaims{
		"sub":   uuid.New().String(),
		"email": "user@example.com",
		"roles": []string{"user"},
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := tokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = authService.Refresh(expiredToken)
	assert.ErrorIs(t, err, AppErr.ErrInvalidToken)
}
