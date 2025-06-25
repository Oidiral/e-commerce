package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

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
)

func setupRSA(t *testing.T) (*config.Config, *rsa.PrivateKey) {
	t.Helper()
	dir := t.TempDir()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	privPath := filepath.Join(dir, "private.pem")
	if err := os.WriteFile(privPath, privPem, 0600); err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pubPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	pubPath := filepath.Join(dir, "public.pem")
	if err := os.WriteFile(pubPath, pubPem, 0644); err != nil {
		t.Fatalf("failed to write public key: %v", err)
	}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			PrivateKeyPath:  privPath,
			PublicKeyPath:   pubPath,
			KeyID:           "test-key-id",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	}
	return cfg, privKey
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	mockClient := mocks.NewMockClientRepository(ctrl)
	cfg, _ := setupRSA(t)
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

	authService := NewAuthService(mockRepo, logger, cfg, mockClient)

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
	mockClient := mocks.NewMockClientRepository(ctrl)
	cfg, _ := setupRSA(t)
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

	authService := NewAuthService(mockRepo, logger, cfg, mockClient)

	_, err := authService.Login(context.Background(), "test2@example.com", "wrongPassword")

	assert.ErrorIs(t, err, AppErr.ErrInvalidCredentials)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	mockClient := mocks.NewMockClientRepository(ctrl)
	cfg, _ := setupRSA(t)
	logger := zerolog.Nop()

	mockRepo.
		EXPECT().
		GetByEmail(gomock.Any(), gomock.Eq("nouser@example.com")).
		Return(nil, AppErr.ErrNotFound)

	authService := NewAuthService(mockRepo, logger, cfg, mockClient)

	_, err := authService.Login(context.Background(), "nouser@example.com", "anyPassword")

	assert.ErrorIs(t, err, AppErr.ErrInvalidCredentials)
}

func TestAuthService_RegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	mockClient := mocks.NewMockClientRepository(ctrl)
	cfg, _ := setupRSA(t)
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

	authService := NewAuthService(mockRepo, logger, cfg, mockClient)

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
	mockClient := mocks.NewMockClientRepository(ctrl)
	cfg, _ := setupRSA(t)
	logger := zerolog.Nop()

	mockRepo.
		EXPECT().
		CreateIfNotExists(gomock.Any(), "existing@example.com", gomock.Any()).
		Return(nil, AppErr.ErrUserAlreadyExists)

	authService := NewAuthService(mockRepo, logger, cfg, mockClient)

	_, err := authService.RegisterUser(context.Background(), "existing@example.com", "anyPassword")

	assert.ErrorIs(t, err, AppErr.ErrUserAlreadyExists)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	// No mocks needed for repository
	cfg, privKey := setupRSA(t)
	logger := zerolog.Nop()
	authService := NewAuthService(nil, logger, cfg, nil)

	userId := uuid.New()
	claims := jwt.MapClaims{
		"sub":   userId.String(),
		"email": "user@example.com",
		"roles": []string{"user"},
		"exp":   time.Now().Add(cfg.JWT.RefreshTokenTTL).Unix(),
		"iat":   time.Now().Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken, err := tokenObj.SignedString(privKey)
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
	cfg, privKey := setupRSA(t)
	logger := zerolog.Nop()
	authService := NewAuthService(nil, logger, cfg, nil)

	claims := jwt.MapClaims{
		"sub":   uuid.New().String(),
		"email": "user@example.com",
		"roles": []string{"user"},
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	expiredToken, err := tokenObj.SignedString(privKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = authService.Refresh(expiredToken)
	assert.ErrorIs(t, err, AppErr.ErrInvalidToken)
}
