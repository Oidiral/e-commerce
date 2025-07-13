package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
	model "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
	appErr "github.com/oidiral/e-commerce/services/auth-svc/internal/errors"
)

const (
	defaultRoleName   = "user"
	pgUniqueViolation = "23505"
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
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	u, err := toDomainFromGetUserByEmailRow(dbUser)
	if err != nil {
		return nil, fmt.Errorf("convert to domain model: %w", err)
	}

	return &u, nil
}

func (r *Repository) CreateIfNotExists(ctx context.Context, email, hash string) (*model.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	qtx := r.q.WithTx(tx)

	dbUser, err := qtx.CreateUserIfNotExists(ctx, db.CreateUserIfNotExistsParams{
		Email:        email,
		PasswordHash: hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, appErr.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	role, err := qtx.GetRoleByName(ctx, defaultRoleName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("default role '%s' not found: %w", defaultRoleName, appErr.ErrRoleNotFound)
		}
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	if err = qtx.CreateUserRole(ctx, db.CreateUserRoleParams{
		UserID: dbUser.ID,
		RoleID: role.ID,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, appErr.ErrInvalidCredentials
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	u, err := toDomainFromAuthUserAndRole(dbUser, role)
	if err != nil {
		return nil, fmt.Errorf("convert to domain model: %w", err)
	}

	return &u, nil
}
