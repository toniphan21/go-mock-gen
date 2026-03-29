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

func Test_STUB_Call_Has_Location(t *testing.T) {
	mock := testTarget()

	spy := mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})

	mock.Full(context.Background(), "anything")

	want := "passed_test.go:31" // 2 lines above
	if spy.Calls[0].Location != want {
		t.Fatalf("want: %s, got: %s", want, spy.Calls[0].Location)
	}
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
