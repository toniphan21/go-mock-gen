package integration

import "context"

type User struct {
	ID       string
	Username string
}

type Repository interface {
	Begin()

	Count() int

	CreateUser(ctx context.Context, user User) error

	GetUsers(ctx context.Context, tenantID string) ([]User, error)
}

type Service interface {
	VeryComplicated(id int) error
}
