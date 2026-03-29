package emitter

import (
	"reflect"
	"testing"
)

func repositoryCompareByReflectEqual[M repositoryMockMethod, T any](m M, argName string, want T, got T, tb testing.TB, expectAt string, index int) {
	if reflect.DeepEqual(want, got) {
		return
	}

	tb.Helper()
	m.fatal(index, repositoryMessageArgumentMismatched(m, argName, expectAt, "reflect.DeepEqual", index+1, want, got))
}
