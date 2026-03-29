package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageMatchFail(m repositoryMockMethod, matchedAt string, index int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s call #%d did not match\n", m.interfaceName(), m.methodName(), index+1))
	sb.WriteString(fmt.Sprintf("arguments:\n"))
	repositoryMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	m.buildCallHistory(sb, "call history")
	sb.WriteString(fmt.Sprintf("hint: check the callback passed to Match at %s", matchedAt))
	return sb.String()
}
