package meta

import (
	"context"
	"fmt"
	"testing"
)

// ExampleData.GenerateCode()
func Test_Example_Target_Full(t *testing.T) {
	mock := testTarget()

	// done - ExampleData.expectCalledCode()
	t.Run("expect called once", func(t *testing.T) {
		mock.EXPECT().Full(t)

		var ctx context.Context
		var input string
		mock.Full(ctx, input)
	})

	// done - ExampleData.expectCalledCode()
	t.Run("expect called twice", func(t *testing.T) {
		mock.EXPECT().Full(t)
		mock.EXPECT().Full(t)

		var ctx context.Context
		var input string
		mock.Full(ctx, input)
		mock.Full(ctx, input)
	})

	// done - ExampleData.expectCalledStubReturnCode()
	t.Run("expect called - stub return", func(t *testing.T) {
		var first []Result
		var second error

		mock.EXPECT().Full(t).Return(first, second)

		var ctx context.Context
		var input string
		ret0, ret1 := mock.Full(ctx, input)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	// done - ExampleData.expectAllUseValue()
	t.Run("expect called - match all arguments by values", func(t *testing.T) {
		var ctx context.Context
		var input string

		mock.EXPECT().Full(t).With(ctx, input)

		mock.Full(ctx, input)
	})

	// done - ExampleData.expectAllUseValueStubReturn()
	t.Run("expect called - match all arguments by values - stub return", func(t *testing.T) {
		var ctx context.Context
		var input string
		var first []Result
		var second error

		mock.EXPECT().Full(t).With(ctx, input).Return(first, second)

		ret0, ret1 := mock.Full(ctx, input)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	// done - ExampleData.expectPartialUseValue()
	t.Run("expect called - match partial argument by value", func(t *testing.T) {
		var ctx context.Context

		mock.EXPECT().Full(t).WithCtx(ctx)

		var input string
		mock.Full(ctx, input)
	})

	// done - ExampleData.expectPartialUseValueStubReturn()
	t.Run("expect called - match partial argument by value - stub result", func(t *testing.T) {
		var ctx context.Context
		var first []Result
		var second error

		mock.EXPECT().Full(t).WithCtx(ctx).Return(first, second)

		var input string
		ret0, ret1 := mock.Full(ctx, input)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	// done - ExampleData.expectAllUseCallback()
	t.Run("expect called - match all arguments by callback", func(t *testing.T) {
		mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
			return true
		})

		var ctx context.Context
		var input string
		mock.Full(ctx, input)
	})

	// done - ExampleData.expectAllUseCallbackStubReturn()
	t.Run("expect called - match all arguments by callback - stub return", func(t *testing.T) {
		var first []Result
		var second error

		mock.EXPECT().Full(t).Match(func(ctx context.Context, input string) bool {
			return true
		}).Return(first, second)

		var ctx context.Context
		var input string
		ret0, ret1 := mock.Full(ctx, input)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	// done - ExampleData.expectPartialUseCallback()
	t.Run("expect called - match partial argument by callback", func(t *testing.T) {
		mock.EXPECT().Full(t).MatchCtx(func(ctx context.Context) bool {
			return true
		})

		var ctx context.Context
		var input string
		mock.Full(ctx, input)
	})

	// done - ExampleData.expectPartialUseCallbackStubReturn()
	t.Run("expect called - match partial argument by callback - stub result", func(t *testing.T) {
		var first []Result
		var second error

		mock.EXPECT().Full(t).MatchCtx(func(ctx context.Context) bool {
			return true
		}).Return(first, second)

		var ctx context.Context
		var input string
		ret0, ret1 := mock.Full(ctx, input)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	// ExampleData.stubCode()
	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
			return []Result{}, nil
		})

		var ctx context.Context
		var input string
		mock.Full(ctx, input)

		fmt.Println(spy)
	})
}
