package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
)

type AuthRepository struct {
	q *db.Queries
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{q: db.New(pool)}
}
