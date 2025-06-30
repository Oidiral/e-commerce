package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/oidiral/e-commerce/services/auth-svc/config"
	user "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	AppErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
	repository "github.com/oidiral/e-commerce/services/auth-svc/internal/repository"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/utils"
	"github.com/rs/zerolog"
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

func (s *AuthService) loadKeys() error {
	privPEM, err := os.ReadFile(s.cfg.JWT.PrivateKeyPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return errors.New("failed to decode private key PEM")
	}
	var parsedKey any
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return errors.New("failed to parse private key")
		}
	}
	s.privKey = parsedKey.(*rsa.PrivateKey)
	jwkKey, err := jwk.FromRaw(s.privKey.Public())
	if err != nil {
		return errors.New("failed to create JWK from public key")
	}

	_ = jwkKey.Set(jwk.KeyIDKey, s.cfg.JWT.KeyID)
	_ = jwkKey.Set(jwk.AlgorithmKey, "RS256")
	_ = jwkKey.Set(jwk.KeyUsageKey, "sig")

	keySet := jwk.NewSet()
	keySet.AddKey(jwkKey)

	buf, err := json.Marshal(keySet)
	if err != nil {
		return err
	}
	s.jwks = buf
	s.log.Info().Msg("JWT keys loaded successfully")
	return nil
}
