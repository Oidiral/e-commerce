package user

import (
	"fmt"

	"github.com/google/uuid"
	db "github.com/oidiral/e-commerce/services/auth-svc/db/sqlc"
	domain "github.com/oidiral/e-commerce/services/auth-svc/internal/domain/model"
)

func toDomainFromAuthUserAndRole(u db.AuthUser, role db.AuthRole) (domain.User, error) {
	id, err := uuid.FromBytes(u.ID.Bytes[:])
	if err != nil {
		return domain.User{}, fmt.Errorf("invalid UUID from AuthUser.ID: %w", err)
	}

	return domain.User{
		ID:        id,
		Email:     u.Email,
		Password:  u.PasswordHash,
		Status:    u.Status,
		Role:      role.Name,
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
		Role:      row.RoleName,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

// Если в будущем понадобится конвертировать domain.User → DB-параметры (например, для обновления),
// можно добавить функции вида fromDomainToAuthUserParams и т.д.:
//
// func fromDomainToCreateUserParams(u domain.User) db.CreateUserParams {
//     return db.CreateUserParams{
//         ID:           u.ID,
//         Email:        u.Email,
//         PasswordHash: u.Password,
//         Status:       int16(u.Status),
//     }
// }
