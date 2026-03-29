package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageDuplicateStub(m repositoryMockMethod, firstUsedAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("duplicate STUB for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "first used at", firstUsedAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "second used at", repositoryCallerLocation(3)))
	sb.WriteString(fmt.Sprintf("\thint: %s.%s is already stubbed, remove one of the above\n\n", m.interfaceName(), m.methodName()))
	return sb.String()
}
