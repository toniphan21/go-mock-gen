package emitter

import "testing"

func repositoryCompareByBasicComparison[M repositoryMockMethod, T comparable](m M, argName string, want T, got T, tb testing.TB, expectAt string, index int) {
	if want == got {
		return
	}

	tb.Helper()
	m.fatal(index, repositoryMessageArgumentMismatched(m, argName, expectAt, "==", index+1, want, got))
}
