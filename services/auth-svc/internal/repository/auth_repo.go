package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	model "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
	appErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
)

type AuthRepository interface {
	CreateIfNotExists(ctx context.Context, email, hash string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type Repository struct {
	q  *db.Queries
	db *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &Repository{
		q:  db.New(pool),
		db: pool,
	}
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErr.ErrNotFound
		}
		return nil, fmt.Errorf("get model by email: %w", err)
	}
	u, err := toDomainFromGetUserByEmailRow(dbUser)
	if err != nil {
		return nil, fmt.Errorf("to domain from get model by email row: %w", err)
	}
	return &u, nil
}

func (r *Repository) CreateIfNotExists(ctx context.Context, email, hash string) (*model.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return &model.User{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	dbUser, err := qtx.CreateUserIfNotExists(ctx, db.CreateUserIfNotExistsParams{
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		return &model.User{}, appErr.ErrUserAlreadyExists
	}

	role, err := qtx.GetRoleByName(ctx, "model")
	if err != nil {
		return &model.User{}, appErr.ErrRoleNotFound
	}

	if err = qtx.CreateUserRole(ctx, db.CreateUserRoleParams{
		UserID: dbUser.ID,
		RoleID: role.ID,
	}); err != nil {
		return &model.User{}, fmt.Errorf("create model-role: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return &model.User{}, fmt.Errorf("commit: %w", err)
	}

	u, err := toDomainFromAuthUserAndRole(dbUser, role)
	if err != nil {
		return &model.User{}, err
	}
	return &u, nil
}
