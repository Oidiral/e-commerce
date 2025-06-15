package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/model"
	db "github.com/oidiral/e-commerce/services/cart-svc/internal/repository/sqlc"
)

var (
	ErrCartNotFound  = errors.New("cart not found")
	ErrItemNotFound  = errors.New("cart item not found")
	ErrQtyConstraint = errors.New("qty constraint violated")
	ErrDB            = errors.New("db failure")
)

type CartRepository interface {
	Get(ctx context.Context, cartID uuid.UUID) (*model.Cart, error)
	ChangeQuantity(ctx context.Context, cartID, productID uuid.UUID, qty int) error
	UpsertItem(ctx context.Context, cartID, productID uuid.UUID, price float64, qty int) error
	DeleteItem(ctx context.Context, cartID, productID uuid.UUID) error
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
}

type cartRepoPg struct {
	q  *db.Queries
	db *pgxpool.Pool
}

func NewCartRepoPg(pool *pgxpool.Pool) CartRepository {
	return &cartRepoPg{
		q:  db.New(pool),
		db: pool,
	}
}

func (r *cartRepoPg) Get(ctx context.Context, cartID uuid.UUID) (*model.Cart, error) {
	dbCart, err := r.q.GetCart(ctx, cartID)
	if err != nil {
		return nil, mapPgErr(err)
	}
	dbItems, err := r.q.ListItems(ctx, cartID)
	if err != nil {
		return nil, mapPgErr(err)
	}
	return mapCartToDomain(dbCart, dbItems), nil
}

func (r *cartRepoPg) UpsertItem(ctx context.Context, cartID, productID uuid.UUID, price float64, qty int) error {
	if qty <= 0 {
		return ErrQtyConstraint
	}

	err := r.q.UpsertCartItem(ctx, db.UpsertCartItemParams{
		CartID:    cartID,
		ProductID: productID,
		Price:     price,
		Quantity:  int32(qty),
	})
	return mapPgErr(err)
}

func (r *cartRepoPg) ChangeQuantity(ctx context.Context, cartID, productID uuid.UUID, qty int) error {
	if qty <= 0 {
		return ErrQtyConstraint
	}
	tag, err := r.q.UpdateQuantity(ctx, db.UpdateQuantityParams{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  int32(qty),
	})
	if err != nil {
		return mapPgErr(err)
	}
	if tag == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *cartRepoPg) DeleteItem(ctx context.Context, cartID, productID uuid.UUID) error {
	tag, err := r.q.DeleteCartItem(ctx, db.DeleteCartItemParams{
		CartID:    cartID,
		ProductID: productID,
	})
	if err != nil {
		return mapPgErr(err)
	}
	if tag == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *cartRepoPg) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	tag, err := r.q.DeleteCart(ctx, cartID)
	if err != nil {
		return mapPgErr(err)
	}
	if tag == 0 {
		return ErrCartNotFound
	}
	return nil
}

func mapCartToDomain(c db.Cart, items []db.CartItem) *model.Cart {
	domainItems := make([]model.CartItem, 0, len(items))
	for _, it := range items {
		domainItems = append(domainItems, model.CartItem{
			ProductID: it.ProductID,
			Price:     it.Price,
			Qty:       int(it.Quantity),
		})
	}
	return &model.Cart{
		ID:        c.ID,
		UserID:    c.UserID,
		Status:    model.CartStatus(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Items:     domainItems,
	}
}

func mapPgErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return ErrCartNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDB
		case "23514":
			return ErrQtyConstraint
		case "23503":
			return ErrCartNotFound
		default:
			return ErrDB
		}
	}

	return err
}
