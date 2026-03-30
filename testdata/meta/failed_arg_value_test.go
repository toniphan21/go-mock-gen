package meta

import (
	"context"
	"testing"
)

func Test_Arg_Value_Duplicate_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithCtx(context.Background()).
		WithCtx(context.Background())

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Value_Duplicate_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithCtx(context.Background()).
		WithInput("a").
		WithInput("b")

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Value_Failed_FirstPlace_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithCtx(context.Background())

	mock.Full(context.WithValue(context.Background(), "key", "val"), "anything")
}

func Test_Arg_Value_Failed_SecondPlace_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithCtx(context.Background()).
		WithInput("a")

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Value_Failed_FirstPlace_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithInput("a")

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Value_Failed_SecondPlace_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		WithCtx(context.Background()).
		WithInput("a")

	mock.Full(context.Background(), "anything")
}
