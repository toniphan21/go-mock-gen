package meta

import (
	"context"
	"errors"
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

	want := "passed_test.go:32" // 2 lines above
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

func Test_EXPECT_Partial_Arg(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t).MatchInput(func(input string) bool {
		return input == "anything"
	})

	mock.Full(context.Background(), "anything")
}

func Test_EXPECT_Return(t *testing.T) {
	mock := testTarget()

	assertOk := func(val string) {
		_, err := mock.Full(context.Background(), "anything")
		if err.Error() != val {
			t.Fatalf("want: %s, got: %s", val, err.Error())
		}
	}

	mock.EXPECT().Full(t).Return(nil, errors.New("whatever-1"))
	assertOk("whatever-1")

	mock.EXPECT().Full(t).With(context.Background(), "anything").
		Return(nil, errors.New("whatever-2"))
	assertOk("whatever-2")

	mock.EXPECT().Full(t).WithInput("anything").
		Return(nil, errors.New("whatever-3"))
	assertOk("whatever-3")

	mock.EXPECT().Full(t).
		Match(func(ctx context.Context, input string) bool { return input == "anything" }).
		Return(nil, errors.New("whatever-4"))
	assertOk("whatever-4")

	mock.EXPECT().Full(t).
		MatchInput(func(input string) bool { return input == "anything" }).
		Return(nil, errors.New("whatever-5"))
	assertOk("whatever-5")
}

func Test_SubTests(t *testing.T) {
	mock := testTarget()

	t.Run("case 1", func(t *testing.T) {
		mock.EXPECT().Full(t)

		mock.Full(context.Background(), "anything")
	})

	t.Run("case 2", func(t *testing.T) {
		mock.EXPECT().Full(t).With(context.Background(), "anything")

		mock.Full(context.Background(), "anything")
	})

	t.Run("case 3", func(t *testing.T) {
		mock.EXPECT().Full(t).MatchInput(func(input string) bool {
			return input == "anything"
		})

		mock.Full(context.Background(), "anything")
	})
}
