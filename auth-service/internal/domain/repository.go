package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateVerified(ctx context.Context, userID string, verified bool) error
}