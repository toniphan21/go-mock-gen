package meta

import (
	"context"
	"testing"
)

func Test_Two_EXPECT_Call_Once(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "a")
}

func Test_Three_EXPECT_Call_Twice(t *testing.T) {
	mock := testTarget()

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	mock.Full(context.Background(), "a")
	mock.Full(context.Background(), "b")
}

func Test_Two_EXPECT_Call_Once_In_Production(t *testing.T) {
	mock := testTarget()
	prod := Production{target: mock}

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	prod.CallFullOnce(context.Background(), "any")
}

func Test_Three_EXPECT_Call_Twice_In_Production(t *testing.T) {
	mock := testTarget()
	prod := Production{target: mock}

	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)
	mock.EXPECT().Full(t)

	prod.CallFullTwice(context.Background(), "any")
}
