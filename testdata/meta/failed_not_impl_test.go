package meta

import (
	"context"
	"testing"
)

func Test_Not_Implemented(t *testing.T) {
	mock := testTarget()

	_, _ = mock.Full(context.Background(), "anything")
}

func Test_Not_Implemented_Via_Ctor(t *testing.T) {
	mock := &target{}

	_, _ = mock.Full(context.Background(), "anything")
}

func Test_Not_Implemented_Via_Another_Func(t *testing.T) {
	mock := testTarget()

	testNotImplementedViaAnotherFunc(mock)
}

func testNotImplementedViaAnotherFunc(mock Target) {
	_, _ = mock.Full(context.Background(), "anything")
}
