package meta

import (
	"context"
	"testing"
)

func Test_CallWith_Fail_FirstCall_FirstArgument(t *testing.T) {
	mock := testTarget()
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "a")

	mock.Full(context.WithValue(ctx, "key", "val"), "a")
}

func Test_CallWith_Fail_FirstCall_SecondArgument(t *testing.T) {
	mock := testTarget()
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "a")

	mock.Full(ctx, "1")
}

func Test_CallWith_Fail_SecondCall_FirstArgument(t *testing.T) {
	mock := testTarget()
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "1")
	mock.EXPECT().Full(t).With(ctx, "1")

	mock.Full(ctx, "1")
	mock.Full(context.WithValue(ctx, "key", "val"), "a")
}

func Test_CallWith_Fail_SecondCall_SecondArgument(t *testing.T) {
	mock := testTarget()
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "1")
	mock.EXPECT().Full(t).With(ctx, "1")

	mock.Full(context.Background(), "1")
	mock.Full(context.Background(), "a")
}

func Test_CallWith_Fail_FirstCall_FirstArgument_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "a")

	prod.CallFullOnce(context.WithValue(ctx, "key", "val"), "a")
}

func Test_CallWith_Fail_FirstCall_SecondArgument_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "a")

	prod.CallFullOnce(ctx, "1")
}

func Test_CallWith_Fail_SecondCall_FirstArgument_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "1")
	mock.EXPECT().Full(t).With(ctx, "1")

	mock.Full(ctx, "1")
	prod.CallFullOnce(context.WithValue(ctx, "key", "val"), "1")
}

func Test_CallWith_Fail_SecondCall_SecondArgument_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}
	ctx := context.Background()

	mock.EXPECT().Full(t).With(ctx, "1")
	mock.EXPECT().Full(t).With(ctx, "1")

	mock.Full(ctx, "1")
	prod.CallFullOnce(context.Background(), "a")
}
