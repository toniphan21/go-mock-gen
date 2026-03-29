package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageStubAfterExpect(m repositoryMockMethod, expectAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "EXPECT used at", expectAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "STUB used at", repositoryCallerLocation(3)))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}
