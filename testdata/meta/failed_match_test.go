package meta

import (
	"context"
	"testing"
)

func Test_Match_Fail_FirstCall(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "1"
	})

	mock.Full(context.Background(), "a")
}

func Test_Match_Fail_SecondCall(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "1"
	})
	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "1"
	})

	mock.Full(context.Background(), "1")
	mock.Full(context.Background(), "a")
}

func Test_Match_Fail_FirstCall_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}

	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "1"
	})

	prod.CallFullOnce(context.Background(), "a")
}

func Test_Match_Fail_SecondCall_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}

	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "a 1"
	})
	mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
		return input == "1"
	})

	prod.CallFullTwice(context.Background(), "a")
}
