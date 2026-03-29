package meta

import (
	"context"
	"testing"
)

func Test_STUB_Via_Struct(t *testing.T) {
	mock := &target{}

	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
}

func Test_STUB_Via_Ctor(t *testing.T) {
	mock := testTarget()

	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
}

func Test_EXPECT_Via_Struct(t *testing.T) {
	mock := &target{}

	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "anything")
}

func Test_EXPECT_Via_Ctor(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "anything")
}
