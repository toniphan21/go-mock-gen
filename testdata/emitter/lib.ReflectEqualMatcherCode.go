package emitter

import "reflect"

func repositoryReflectEqualMatcher[T any](want T) func(T) bool {
	return func(got T) bool {
		return reflect.DeepEqual(want, got)
	}
}
