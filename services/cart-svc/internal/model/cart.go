package model

import (
	"github.com/google/uuid"
	"time"
)

type CartStatus string

const (
	CartOpen      CartStatus = "OPEN"
	CartPending   CartStatus = "PENDING"
	CartCheckout  CartStatus = "CHECKOUT"
	CartAbandoned CartStatus = "ABANDONED"
)

type CartItem struct {
	CartID    uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
	Qty       int       `json:"qty"`
}

type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Status    CartStatus `json:"status"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
