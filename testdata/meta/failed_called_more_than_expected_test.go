package meta

import (
	"context"
	"testing"
)

func Test_One_EXPECT_Call_Twice(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "a")
	mock.Full(context.Background(), "b")
}

func Test_Two_EXPECT_Call_Thrice(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "a")
	mock.Full(context.Background(), "b")
	mock.Full(context.Background(), "c")
}

func Test_One_EXPECT_Call_Twice_In_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}

	mock.EXPECT().Full(t)

	prod.CallFullTwice(context.Background(), "any")
}

func Test_Two_EXPECT_Call_Thrice_In_Production(t *testing.T) {
	mock := testTarget()
	prod := &Production{target: mock}

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	prod.CallFullThrice(context.Background(), "any")
}
