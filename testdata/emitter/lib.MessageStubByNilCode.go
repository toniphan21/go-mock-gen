package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageStubByNil(m repositoryMockMethod, calledAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s STUB received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("called at: %s\n\n", calledAt))
	sb.WriteString("hint: provide a valid function\n")
	return sb.String()
}
