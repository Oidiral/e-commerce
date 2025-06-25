package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	user "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	AppErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	repository "github.com/oidiral/e-commerce/services/auth-svc/internal/repository"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/utils"
	"github.com/rs/zerolog"
	"math/big"
	"os"
	"time"
)

type AuthService struct {
	repo    repository.AuthRepository
	log     zerolog.Logger
	cfg     *config.Config
	cliRepo repository.ClientRepository
	privKey *rsa.PrivateKey
	jwks    json.RawMessage
}

func NewAuthService(repo repository.AuthRepository, log zerolog.Logger, cfg *config.Config, cliRepo repository.ClientRepository) *AuthService {
	svc := &AuthService{repo: repo, log: log, cfg: cfg, cliRepo: cliRepo}
	svc.loadKeys()
	return svc
}

type TokenPair struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	AccessExpiresAt  int64  `json:"accessExpiresAt"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}
type TokenService struct {
	AccessToken     string `json:"accessToken"`
	AccessExpiresAt int64  `json:"accessExpiresAt"`
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
		return &s.privKey.PublicKey, nil
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
	rolesClaim, ok := claims["roles"].([]interface{})
	if !ok {
		s.log.Error().Msg("invalid token roles")
		return nil, AppErr.ErrInvalidToken
	}
	roles := make([]string, len(rolesClaim))
	for i, v := range rolesClaim {
		r, ok := v.(string)
		if !ok {
			s.log.Error().Msg("invalid token roles")
			return nil, AppErr.ErrInvalidToken
		}
		roles[i] = r
	}
	u := user.User{
		ID:    subUUID,
		Email: email,
		Roles: roles,
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
		"roles": u.Roles,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.cfg.JWT.KeyID
	jw, err := token.SignedString(s.privKey)
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
		"roles": u.Roles,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.cfg.JWT.KeyID
	jw, err := token.SignedString(s.privKey)
	if err != nil {
		s.log.Error().Err(err).Msg("create refresh token")
		return "", time.Time{}, err
	}
	return jw, exp, nil
}

func (s *AuthService) ClientToken(ctx context.Context, id string, secret string) (*TokenService, error) {
	cli, err := s.cliRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, AppErr.ErrNotFound) {
			return nil, AppErr.ErrInvalidCredentials
		}
		s.log.Error().Err(err).Msg("get client by id")
		return nil, err
	}
	if !utils.CheckPasswordHash(secret, cli.Secret) {
		return nil, AppErr.ErrInvalidCredentials
	}
	now := time.Now()
	exp := now.Add(s.cfg.JWT.AccessTokenTTL)
	claims := jwt.MapClaims{
		"sub":   cli.ID,
		"roles": cli.Roles,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.cfg.JWT.KeyID
	jw, err := token.SignedString(s.privKey)
	if err != nil {
		s.log.Error().Err(err).Msg("create client token")
		return nil, err
	}
	return &TokenService{
		AccessToken:     jw,
		AccessExpiresAt: exp.Unix(),
	}, nil
}

func (s *AuthService) JWKS() []byte {
	return s.jwks
}

func (s *AuthService) loadKeys() {
	privBytes, err := os.ReadFile(s.cfg.JWT.PrivateKeyPath)
	if err != nil {
		s.log.Fatal().Err(err).Msg("read private key")
	}
	block, _ := pem.Decode(privBytes)
	if block == nil {
		s.log.Fatal().Msg("failed to parse PEM block containing the private key")
	}
	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		pk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			s.log.Fatal().Err(err).Msg("parse private key")
		}
		s.privKey = pk.(*rsa.PrivateKey)
	} else {
		s.privKey = pk.(*rsa.PrivateKey)
	}
	pubBytes, err := os.ReadFile(s.cfg.JWT.PublicKeyPath)
	if err != nil {
		s.log.Fatal().Err(err).Msg("read public key")
	}
	pblock, _ := pem.Decode(pubBytes)
	if pblock == nil {
		s.log.Fatal().Msg("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(pblock.Bytes)
	if err != nil {
		s.log.Fatal().Err(err).Msg("parse public key")
	}
	rsaPub := pub.(*rsa.PublicKey)

	n := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPub.E)).Bytes())

	jwk := map[string]string{
		"kty": "RSA",
		"alg": "RS256",
		"use": "sig",
		"kid": s.cfg.JWT.KeyID,
		"n":   n,
		"e":   e,
	}

	set := map[string]interface{}{"keys": []interface{}{jwk}}
	data, err := json.Marshal(set)
	if err != nil {
		s.log.Fatal().Err(err).Msg("marshal JWKS")
	}
	s.jwks = data
}
