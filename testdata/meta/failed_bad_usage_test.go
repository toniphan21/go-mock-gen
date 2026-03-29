package meta

import (
	"context"
	"testing"
)

func Test_Use_STUB_Twice(t *testing.T) {
	mock := testTarget()

	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
}

func Test_Use_STUB_Thrice(t *testing.T) {
	mock := testTarget()

	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
	t.Log("Test_Use_STUB_Thrice")
	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
	t.Log("Test_Use_STUB_Thrice")
	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
}

func Test_Use_STUB_After_EXPECT(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
}

func Test_Use_EXPECT_After_STUB(t *testing.T) {
	mock := testTarget()

	mock.STUB().Full(func(ctx context.Context, input string) ([]Result, error) {
		return nil, nil
	})
	mock.EXPECT().Full(t)
}

func Test_Pass_Nil_To_EXPECT(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(nil)
}

func Test_Pass_Nil_To_STUB(t *testing.T) {
	mock := testTarget()

	mock.STUB().Full(nil)
}

func Test_Pass_Nil_To_STUB_After_Expect(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.STUB().Full(nil)
}

func Test_Pass_Nil_To_EXPECT_Match(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).Match(nil)
}

func Test_Pass_Nil_To_EXPECT_After_EXPECT(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(nil)
}
