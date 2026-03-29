package emitter

import "strings"

type repositoryMockMethod interface {
	methodName() string
	interfaceName() string
	buildCallHistory(sb *strings.Builder, header string)
	fatal(index int, msg string)
	panic(msg string)
}
