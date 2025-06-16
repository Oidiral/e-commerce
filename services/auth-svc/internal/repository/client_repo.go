package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	AppErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
)

type ClientRepository interface {
	GetById(ctx context.Context, id string) (*model.Client, error)
}

type CliRepository struct {
	q  *db.Queries
	db *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) ClientRepository {
	return &CliRepository{
		q:  db.New(pool),
		db: pool,
	}
}

func (r CliRepository) GetById(ctx context.Context, id string) (*model.Client, error) {
	client, err := r.q.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, AppErr.ErrNotFound
		}
		return nil, err
	}
	c := toDomainFromGetClientById(client)
	return &c, nil
}
