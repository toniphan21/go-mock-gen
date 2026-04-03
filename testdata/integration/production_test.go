package integration

import (
	"context"
	"testing"
)

func Test_createUserOp(t *testing.T) {
	repo := testRepository()

	user := User{ID: "1", Username: "user"}
	repo.EXPECT().CreateUser(t).WithUser(User{ID: "1", Username: "user"})

	op := createUserOp{repo: repo}

	err := op.Handle(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
}
