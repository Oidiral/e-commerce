package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/db"
	"time"

	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/model"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
)

const (
	cacheKeyPrefix = "cart:"
	cacheTTL       = 30 * 24 * time.Hour
	bgOpTimeout    = 5 * time.Second
)

type CartService interface {
	GetCart(ctx context.Context, cartID uuid.UUID) (*model.Cart, error)
	AddItem(ctx context.Context, cartID, productID uuid.UUID, price float64, qty int) error
	ChangeQty(ctx context.Context, cartID, productID uuid.UUID, qty int) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	Clear(ctx context.Context, cartID uuid.UUID) error
	Checkout(ctx context.Context, cartID uuid.UUID, paymentMethodID string) error
}

type CartSvc struct {
	db  postgres.CartRepository
	log zerolog.Logger
	rds *redis.Client
}

func NewCartService(dbRepo postgres.CartRepository, logger zerolog.Logger, rdb *db.RedisClient) *CartSvc {
	return &CartSvc{db: dbRepo, log: logger, rds: rdb.Client}
}

func (s *CartSvc) GetCart(ctx context.Context, cartID uuid.UUID) (*model.Cart, error) {
	key := s.cacheKey(cartID)
	raw, err := s.rds.Get(ctx, key).Result()
	if err == nil {
		var cart model.Cart
		if json.Unmarshal([]byte(raw), &cart) == nil {
			s.log.Info().Str("cart_id", cartID.String()).Msg("cache hit for cart")
			return &cart, nil
		}
		s.log.Warn().Str("cart_id", cartID.String()).Msg("failed to unmarshal cached cart, fallback to DB")
	} else if !errors.Is(err, redis.Nil) {
		s.log.Warn().Err(err).Str("cart_id", cartID.String()).Msg("redis GET error, fallback to DB")
	}
	cart, err := s.db.Get(ctx, cartID)
	switch {
	case errors.Is(err, postgres.ErrCartNotFound):
		s.log.Warn().Str("cart_id", cartID.String()).Msg("cart not found")
		return nil, ErrNotFound
	case err != nil:
		s.log.Error().Err(err).Str("cart_id", cartID.String()).Msg("db GetCart failed")
		return nil, ErrInternal
	}
	s.log.Info().Str("cart_id", cartID.String()).Msg("cart loaded from DB")
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		s.refreshCache(bgCtx, cartID, cart)
	}()
	return cart, nil
}

func (s *CartSvc) AddItem(ctx context.Context, cartID, productID uuid.UUID, price float64, qty int) error {
	if qty <= 0 {
		return ErrBadRequest
	}
	if err := s.db.UpsertItem(ctx, cartID, productID, price, qty); err != nil {
		switch {
		case errors.Is(err, postgres.ErrCartNotFound):
			return ErrNotFound
		case errors.Is(err, postgres.ErrQtyConstraint):
			return ErrBadRequest
		default:
			s.log.Error().Err(err).Msg("UpsertItem failed")
			return ErrInternal
		}
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		s.refreshCache(bgCtx, cartID, nil)
	}()
	return nil
}

func (s *CartSvc) ChangeQty(ctx context.Context, cartID, productID uuid.UUID, qty int) error {
	if qty <= 0 {
		return ErrBadRequest
	}
	if err := s.db.ChangeQuantity(ctx, cartID, productID, qty); err != nil {
		switch {
		case errors.Is(err, postgres.ErrItemNotFound):
			return ErrNotFound
		case errors.Is(err, postgres.ErrQtyConstraint):
			return ErrBadRequest
		default:
			s.log.Error().Err(err).Msg("ChangeQty failed")
			return ErrInternal
		}
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		s.refreshCache(bgCtx, cartID, nil)
	}()
	return nil
}

func (s *CartSvc) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	if err := s.db.DeleteItem(ctx, cartID, productID); err != nil {
		switch {
		case errors.Is(err, postgres.ErrItemNotFound):
			return ErrNotFound
		default:
			s.log.Error().Err(err).Msg("DeleteItem failed")
			return ErrInternal
		}
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		s.refreshCache(bgCtx, cartID, nil)
	}()
	return nil
}

func (s *CartSvc) Clear(ctx context.Context, cartID uuid.UUID) error {
	if err := s.db.DeleteCart(ctx, cartID); err != nil {
		switch {
		case errors.Is(err, postgres.ErrCartNotFound):
			return ErrNotFound
		default:
			s.log.Error().Err(err).Msg("DeleteCart failed")
			return ErrInternal
		}
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		if err := s.rds.Del(bgCtx, s.cacheKey(cartID)).Err(); err != nil {
			s.log.Warn().Err(err).Str("cart_id", cartID.String()).Msg("failed to delete cache on Clear")
		}
		s.log.Info().Str("cart_id", cartID.String()).Msg("cart cleared and cache deleted")
	}()
	return nil
}

func (s *CartSvc) Checkout(ctx context.Context, cartID uuid.UUID, paymentMethodID string) error {
	panic("implement me")
}

func (s *CartSvc) cacheKey(cartID uuid.UUID) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, cartID.String())
}

func (s *CartSvc) refreshCache(ctx context.Context, cartID uuid.UUID, cartVal *model.Cart) {
	var cart *model.Cart
	var err error
	if cartVal != nil {
		cart = cartVal
	} else {
		cart, err = s.db.Get(ctx, cartID)
		if err != nil {
			s.log.Warn().Err(err).Str("cart_id", cartID.String()).Msg("cannot refresh cache")
			return
		}
	}
	data, err := json.Marshal(cart)
	if err != nil {
		s.log.Error().Err(err).Msg("failed marshal cart in refreshCache")
		return
	}
	if err := s.rds.Set(ctx, s.cacheKey(cartID), data, cacheTTL).Err(); err != nil {
		s.log.Warn().Err(err).Msg("failed set cache in refreshCache")
	}
}
