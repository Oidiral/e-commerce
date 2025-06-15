package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/model"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/repository/postgres"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
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
}

func NewCartService(db postgres.CartRepository, log zerolog.Logger) *CartSvc {
	return &CartSvc{
		db:  db,
		log: log,
	}
}

func (s *CartSvc) GetCart(ctx context.Context, cartID uuid.UUID) (*model.Cart, error) {
	cart, err := s.db.Get(ctx, cartID)
	switch {
	case errors.Is(err, postgres.ErrCartNotFound):
		s.log.Warn().Msgf("cart with ID %s not found", cartID)
		return nil, ErrNotFound
	case err != nil:
		s.log.Warn().Msgf("Iternal error while getting cart with ID %s: %v", cartID, err)
		return nil, ErrInternal
	}
	s.log.Info().Msgf("cart with ID %s retrieved successfully", cartID)
	return cart, nil
}

func (s *CartSvc) AddItem(ctx context.Context, cartID, productID uuid.UUID, price float64, qty int) error {
	err := s.db.UpsertItem(ctx, cartID, productID, price, qty)
	switch {
	case errors.Is(err, postgres.ErrCartNotFound):
		s.log.Warn().Msgf("cart with ID %s not found", cartID)
		return ErrNotFound
	case errors.Is(err, postgres.ErrQtyConstraint):
		s.log.Warn().Msgf("quantity constraint violated for product %s in cart %s", productID, cartID)
		return ErrBadRequest
	case err != nil:
		s.log.Error().Msgf("internal error while adding item with product ID %s to cart with ID %s: %v", productID, cartID, err)
		return ErrInternal
	}
	s.log.Info().Msgf("item with product ID %s added to cart with ID %s successfully", productID, cartID)
	return nil
}

func (s *CartSvc) ChangeQty(ctx context.Context, cartID, productID uuid.UUID, qty int) error {
	err := s.db.ChangeQuantity(ctx, cartID, productID, qty)
	switch {
	case errors.Is(err, postgres.ErrItemNotFound):
		s.log.Warn().Msgf("cart with ID %s not found", cartID)
		return ErrNotFound
	case errors.Is(err, postgres.ErrQtyConstraint):
		s.log.Warn().Msgf("quantity constraint violated for product %s in cart %s", productID, cartID)
		return ErrBadRequest
	case err != nil:
		s.log.Error().Msgf("internal error while changing quantity for product ID %s in cart with ID %s: %v", productID, cartID, err)
		return ErrInternal
	}
	s.log.Info().Msgf("quantity for product ID %s in cart with ID %s changed successfully", productID, cartID)
	return nil
}

func (s *CartSvc) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	err := s.db.DeleteItem(ctx, cartID, productID)
	switch {
	case errors.Is(err, postgres.ErrItemNotFound):
		s.log.Warn().Msgf("item with product ID %s not found in cart with ID %s", productID, cartID)
		return ErrNotFound
	case err != nil:
		s.log.Error().Msgf("internal error while removing item with product ID %s from cart with ID %s: %v", productID, cartID, err)
		return ErrInternal
	}
	s.log.Info().Msgf("item with product ID %s removed from cart with ID %s successfully", productID, cartID)
	return nil
}

func (s *CartSvc) Clear(ctx context.Context, cartID uuid.UUID) error {
	err := s.db.DeleteCart(ctx, cartID)
	switch {
	case errors.Is(err, postgres.ErrCartNotFound):
		s.log.Warn().Msgf("cart with ID %s not found", cartID)
		return ErrNotFound
	case err != nil:
		s.log.Error().Msgf("internal error while clearing cart with ID %s: %v", cartID, err)
		return ErrInternal
	}
	s.log.Info().Msgf("cart with ID %s cleared successfully", cartID)
	return nil
}

func (s *CartSvc) Checkout(ctx context.Context, cartID uuid.UUID, paymentMethodID string) error {
	//TODO implement me
	panic("implement me")
}
