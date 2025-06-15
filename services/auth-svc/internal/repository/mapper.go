package user

import (
	"fmt"

	"github.com/google/uuid"
	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
	domain "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
)

func toDomainFromAuthUserAndRole(u db.User, role db.Role) (domain.User, error) {
	id, err := uuid.FromBytes(u.ID.Bytes[:])
	if err != nil {
		return domain.User{}, fmt.Errorf("invalid UUID from AuthUser.ID: %w", err)
	}

	return domain.User{
		ID:        id,
		Email:     u.Email,
		Password:  u.PasswordHash,
		Status:    u.Status,
		Roles:     []string{role.Name},
		CreatedAt: u.CreatedAt.Time,
	}, nil
}

func toDomainFromGetUserByEmailRow(row db.GetUserByEmailRow) (domain.User, error) {
	id, err := uuid.FromBytes(row.ID.Bytes[:])
	if err != nil {
		return domain.User{}, fmt.Errorf("invalid UUID from GetUserByEmailRow.ID: %w", err)
	}

	return domain.User{
		ID:        id,
		Email:     row.Email,
		Password:  row.PasswordHash,
		Status:    row.Status,
		Roles:     row.Roles,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}
