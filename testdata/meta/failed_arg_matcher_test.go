package meta

import (
	"context"
	"testing"
)

func Test_Arg_Matcher_Pass_Nil_In_First_Place_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).MatchCtx(nil)
}

func Test_Arg_Matcher_Pass_Nil_In_Second_Place_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).MatchInput(func(input string) bool {
		return true
	}).MatchCtx(nil)
}

func Test_Arg_Matcher_Pass_Nil_In_First_Place_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).MatchInput(nil)
}

func Test_Arg_Matcher_Pass_Nil_In_Second_Place_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).MatchCtx(func(ctx context.Context) bool {
		return true
	}).MatchInput(nil)
}

func Test_Arg_Matcher_Duplicate_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchCtx(func(ctx context.Context) bool { return true }).
		MatchCtx(func(ctx context.Context) bool { return true })

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Matcher_Duplicate_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchCtx(func(ctx context.Context) bool { return true }).
		MatchInput(func(input string) bool { return false }).
		MatchInput(func(input string) bool { return false })

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Matcher_Failed_FirstPlace_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchCtx(func(ctx context.Context) bool { return false })

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Matcher_Failed_SecondPlace_FirstArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchInput(func(input string) bool { return true }).
		MatchCtx(func(ctx context.Context) bool { return false })

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Matcher_Failed_FirstPlace_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchInput(func(input string) bool { return false })

	mock.Full(context.Background(), "anything")
}

func Test_Arg_Matcher_Failed_SecondPlace_SecondArg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).
		MatchCtx(func(ctx context.Context) bool { return true }).
		MatchInput(func(input string) bool { return false })

	mock.Full(context.Background(), "anything")
}
