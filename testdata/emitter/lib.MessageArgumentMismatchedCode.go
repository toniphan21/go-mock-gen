package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageArgumentMismatched(m repositoryMockMethod, argName string, expectAt string, comparedBy string, callNo int, want any, got any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s call #%d argument \"%s\" did not match\n", m.interfaceName(), m.methodName(), callNo, argName))
	sb.WriteString(fmt.Sprintf("  want: %#v\n", want))
	sb.WriteString(fmt.Sprintf("   got: %#v\n", got))
	sb.WriteString(fmt.Sprintf("method: %s\n", comparedBy))
	sb.WriteString("\n")
	m.buildCallHistory(sb, "call history")
	sb.WriteString(fmt.Sprintf("hint: for custom matching use .Match(func(...) bool) at %s\n\tor use STUB for fine-grained control", expectAt))
	return sb.String()
}
