package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageMatchArgByNil(m repositoryMockMethod, method string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s %s received a nil function\n", m.interfaceName(), m.methodName(), method))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}
