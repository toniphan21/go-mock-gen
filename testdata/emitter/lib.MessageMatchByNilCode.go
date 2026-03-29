package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageMatchByNil(m repositoryMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s Match received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}
