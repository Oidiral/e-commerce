// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserIfNotExists(ctx context.Context, arg CreateUserIfNotExistsParams) (User, error)
	CreateUserRole(ctx context.Context, arg CreateUserRoleParams) error
	GetById(ctx context.Context, id string) (Client, error)
	GetRoleByName(ctx context.Context, name string) (Role, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
}

var _ Querier = (*Queries)(nil)
