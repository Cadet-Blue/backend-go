package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}
