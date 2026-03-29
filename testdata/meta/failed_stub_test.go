package meta

import (
	"context"
	"testing"
)

func Test_Stub_Fail_Not_Called(t *testing.T) {
	mock := testTarget()

	spy := mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})

	if len(spy.Calls) != 1 {
		t.Fatalf("want 1, got %v", len(spy.Calls))
	}
}

func Test_Stub_Fail_Too_Many_Calls(t *testing.T) {
	mock := testTarget()

	spy := mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})

	mock.Full(context.Background(), "a")
	mock.Full(context.Background(), "b")

	if len(spy.Calls) != 1 {
		t.Fatalf("want 1, got %v", len(spy.Calls))
	}
}
