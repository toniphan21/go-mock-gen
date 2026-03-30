package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageDuplicateMatchArg(m repositoryMockMethod, method string, firstUsedAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("duplicate %s for %s.%s\n", method, m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "first used at", firstUsedAt))
	sb.WriteString("\thint: each argument can only be matched once, remove one of the above")
	return sb.String()
}
