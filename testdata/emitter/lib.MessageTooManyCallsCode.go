package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageTooManyCalls(m repositoryMockMethod, want int, got int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("too many calls to %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", got, "missing"))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", repositoryCallerLocation(4)))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	repositoryMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().%s(t)\n", m.methodName()))
	return sb.String()
}
