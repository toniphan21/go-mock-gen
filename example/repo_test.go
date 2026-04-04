package example_test

import (
	"context"
	"testing"

	"nhatp.com/go/mock-gen/example/mock"
)

func Test_Repository(t *testing.T) {
	repo := mock.NewRepository()
	repo.EXPECT().CreateUser(t).WithAge(10)

	_, _ = repo.CreateUser(context.Background(), "anything", 10)
}
