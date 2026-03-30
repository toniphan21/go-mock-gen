package emitter

func repositoryBasicComparisonMatcher[T comparable](want T) func(T) bool {
	return func(got T) bool {
		return want == got
	}
}
