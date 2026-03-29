package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageExpectByNil(m repositoryMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected nil testing.TB in %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\tcalled at: %s\n\n", repositoryCallerLocation(3)))
	sb.WriteString("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n")
	sb.WriteString(fmt.Sprintf("\t\tspy := [var].STUB().%s(func(...) ...)\n", m.methodName()))
	return sb.String()
}
