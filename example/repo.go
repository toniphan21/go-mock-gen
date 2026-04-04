package example

import (
	"context"
)

/* generate code for Repository in package mock
 *   //go:generate go run ../cmd/go-mock-gen -i Repository -s Repository -o mock/mockgen.go
 */

type User struct{}

//go:generate go run ../cmd/go-mock-gen -i Repository

type Repository interface {
	CreateUser(ctx context.Context, username string, age int) (User, error)

	GetUsers(ctx context.Context) ([]User, error)
}
