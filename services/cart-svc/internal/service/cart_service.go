package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/db"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/pb/catalog"
	"time"

	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/model"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrBadRequest  = errors.New("bad request")
	ErrInternal    = errors.New("internal error")
	ErrOutOfStock  = errors.New("out of stock")
	ErrInvalidItem = errors.New("invalid cart item")
)

const (
	cacheKeyPrefix = "cart:"
	cacheTTL       = 30 * 24 * time.Hour
	bgOpTimeout    = 5 * time.Second
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error)
	AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error
	ChangeQty(ctx context.Context, userID, productID uuid.UUID, qty int) error
	RemoveItem(ctx context.Context, userID, productID uuid.UUID) error
	Clear(ctx context.Context, userID uuid.UUID) error
	Checkout(ctx context.Context, userID uuid.UUID) error
}

type CartSvc struct {
	db            postgres.CartRepository
	log           zerolog.Logger
	rds           *redis.Client
	catalogClient catalog.CatalogClient
}

func NewCartService(dbRepo postgres.CartRepository, logger zerolog.Logger, rdb *db.RedisClient, catClient catalog.CatalogClient) *CartSvc {
	return &CartSvc{db: dbRepo, log: logger, rds: rdb.Client, catalogClient: catClient}
}

func (s *CartSvc) GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	key := s.cacheKey(userID)
	raw, err := s.rds.Get(ctx, key).Result()
	if err == nil {
		var cart model.Cart
		if json.Unmarshal([]byte(raw), &cart) == nil {
			s.log.Info().Str("user_id", userID.String()).Msg("cache hit for cart")
			return &cart, nil
		}
		s.log.Warn().Str("user_id", userID.String()).Msg("failed to unmarshal cached cart, fallback to DB")
	} else if !errors.Is(err, redis.Nil) {
		s.log.Warn().Err(err).Str("user_id", userID.String()).Msg("redis GET error, fallback to DB")
	}
	cart, err := s.db.GetByUser(ctx, userID)
	switch {
	case errors.Is(err, postgres.ErrCartNotFound):
		s.log.Warn().Str("user_id", userID.String()).Msg("cart not found")
		return nil, ErrNotFound
	case err != nil:
		s.log.Error().Err(err).Str("user_id", userID.String()).Msg("db GetCart failed")
		return nil, ErrInternal
	}
	s.log.Info().Str("user_id", userID.String()).Msg("cart loaded from DB")
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), bgOpTimeout)
		defer cancel()
		s.refreshCache(bgCtx, userID, cart)
	}()
	return cart, nil
}

func (s *CartSvc) AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
	if qty <= 0 || productID == uuid.Nil {
		return ErrBadRequest
	}
	resp, err := s.catalogClient.GetPriceWithQty(ctx, &catalog.GetPriceRequest{
		ProductId: productID.String(),
	})
	if err != nil {
		s.log.Error().Err(err).Str("product_id", productID.String()).Msg("GetPriceWithQty failed")
		return ErrInternal
	}
	if resp.AvailableQty < int32(qty) {
		s.log.Warn().Str("product_id", productID.String()).Int("requested_qty", qty).Int32("available_qty", resp.AvailableQty).Msg("requested quantity exceeds available stock")
		return ErrBadRequest
	}
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.db.UpsertItem(ctx, cart.ID, productID, float64(resp.Price), qty); err != nil {
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
		s.refreshCache(bgCtx, userID, nil)
	}()
	return nil
}

func (s *CartSvc) ChangeQty(ctx context.Context, userID, productID uuid.UUID, qty int) error {
	if qty <= 0 || productID == uuid.Nil {
		return ErrBadRequest
	}
	resp, err := s.catalogClient.GetQty(ctx, &catalog.GetQtyRequest{
		ProductId: productID.String(),
	})
	if err != nil {
		s.log.Error().Err(err).Str("product_id", productID.String()).Msg("GetQty failed")
		return ErrInternal
	}
	if resp.AvailableQty < int32(qty) {
		s.log.Warn().Str("product_id", productID.String()).Int("requested_qty", qty).Int32("available_qty", resp.AvailableQty).Msg("requested quantity exceeds available stock")
		return ErrBadRequest
	}
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.db.ChangeQuantity(ctx, cart.ID, productID, qty); err != nil {
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
		s.refreshCache(bgCtx, userID, nil)
	}()
	return nil
}

func (s *CartSvc) RemoveItem(ctx context.Context, userID, productID uuid.UUID) error {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.db.DeleteItem(ctx, cart.ID, productID); err != nil {
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
		s.refreshCache(bgCtx, userID, nil)
	}()
	return nil
}

func (s *CartSvc) Clear(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.db.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrCartNotFound) {
			return ErrNotFound
		}
		s.log.Error().Err(err).Msg("GetByUser failed")
		return ErrInternal
	}
	if err := s.db.DeleteCart(ctx, cart.ID); err != nil {
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
		if err := s.rds.Del(bgCtx, s.cacheKey(userID)).Err(); err != nil {
			s.log.Warn().Err(err).Str("user_id", userID.String()).Msg("failed to delete cache on Clear")
		}
		s.log.Info().Str("user_id", userID.String()).Msg("cart cleared and cache deleted")
	}()
	return nil
}

func (s *CartSvc) Checkout(ctx context.Context, userID uuid.UUID) error {
	key := s.cacheKey(userID)
	raw, rErr := s.rds.GetDel(ctx, key).Result()
	if rErr == nil {
		var cart model.Cart
		if unErr := json.Unmarshal([]byte(raw), &cart); unErr == nil {
			s.log.Info().Str("user_id", userID.String()).
				Msg("cache hit for cart during Checkout")
			if err := s.checkProduct(ctx, &cart); err != nil {
				s.log.Error().Err(err).Str("user_id", userID.String()).
					Msg("product check failed during Checkout")
				return err
			}
			s.log.Info().Str("user_id", userID.String()).
				Msg("checkout completed successfully using cache")
			return nil
		}

		s.log.Warn().Str("user_id", userID.String()).Msg("cache miss for checkout")

	} else if !errors.Is(rErr, redis.Nil) {
		s.log.Warn().Err(rErr).Str("user_id", userID.String()).
			Msg("failed to get cart from cache, falling back to DB")
	} else {
		s.log.Info().Str("user_id", userID.String()).
			Msg("cache miss for cart during Checkout, falling back to DB")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.log.Warn().Str("user_id", userID.String()).
				Msg("cart not found during Checkout")
			return ErrNotFound
		}
		s.log.Error().Err(err).Str("user_id", userID.String()).
			Msg("GetCart failed during Checkout")
		return ErrInternal
	}

	if err := s.checkProduct(ctx, cart); err != nil {
		s.log.Error().Err(err).Str("user_id", userID.String()).
			Msg("product check failed during Checkout")
		return err
	}

	if err := s.Clear(ctx, userID); err != nil && !errors.Is(err, ErrNotFound) {
		s.log.Error().Err(err).Str("user_id", userID.String()).
			Msg("Clear failed during Checkout")
		return ErrInternal
	}

	s.log.Info().Str("user_id", userID.String()).
		Msg("checkout completed successfully using DB")
	return nil
}

func (s *CartSvc) checkProduct(ctx context.Context, cart *model.Cart) error {
	if cart == nil || len(cart.Items) == 0 {
		return ErrNotFound
	}

	for _, item := range cart.Items {
		if item.ProductID == uuid.Nil || item.Qty <= 0 {
			return ErrInvalidItem
		}

		resp, err := s.catalogClient.Checkout(ctx, &catalog.CheckoutRequest{
			ItemId:   item.ProductID.String(),
			Quantity: int32(item.Qty),
		})
		if err != nil {
			s.log.Error().Err(err).Str("product_id", item.ProductID.String()).
				Msg("catalog Checkout RPC failed")
			return ErrInternal
		}

		if !resp.GetAvailable() {
			s.log.Warn().Str("product_id", item.ProductID.String()).
				Msg("product not available for checkout")
			return ErrOutOfStock
		}
	}
	return nil
}

func (s *CartSvc) cacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, userID.String())
}

func (s *CartSvc) refreshCache(ctx context.Context, userID uuid.UUID, cartVal *model.Cart) {
	var cart *model.Cart
	var err error
	if cartVal != nil {
		cart = cartVal
	} else {
		cart, err = s.db.GetByUser(ctx, userID)
		if err != nil {
			s.log.Warn().Err(err).Str("user_id", userID.String()).Msg("cannot refresh cache")
			return
		}
	}
	data, err := json.Marshal(cart)
	if err != nil {
		s.log.Error().Err(err).Msg("failed marshal cart in refreshCache")
		return
	}
	if err := s.rds.Set(ctx, s.cacheKey(userID), data, cacheTTL).Err(); err != nil {
		s.log.Warn().Err(err).Msg("failed set cache in refreshCache")
	}
}

func (s *CartSvc) getOrCreateCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	cart, err := s.db.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrCartNotFound) {
			cart, err = s.db.Create(ctx, userID)
			if err != nil {
				s.log.Error().Err(err).Str("user_id", userID.String()).Msg("create cart failed")
				return nil, ErrInternal
			}
			return cart, nil
		}
		s.log.Error().Err(err).Str("user_id", userID.String()).Msg("GetByUser failed")
		return nil, ErrInternal
	}
	return cart, nil
}
