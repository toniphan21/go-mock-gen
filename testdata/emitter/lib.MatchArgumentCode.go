package emitter

import (
	"fmt"
	"strings"
	"testing"
)

func repositoryMatchArgument[T any](m repositoryMockMethod, index int, name string, got T, match func(T) bool, wants map[string]any, methods map[string]string, hints map[string]string, tb testing.TB, expectAt string) {
	if match == nil || match(got) {
		return
	}
	tb.Helper()

	method := "func(got) bool"
	if v, ok := methods[name]; ok {
		method = v
	}

	hint := fmt.Sprintf("hint: for custom matching use .Match[arg](func(...) bool) at %s\n\tor use STUB for fine-grained control", expectAt)
	if v, ok := hints[name]; ok {
		hint = v
	}

	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s call #%d argument \"%s\" did not match\n", m.interfaceName(), m.methodName(), index+1, name))
	if want, ok := wants[name]; ok {
		sb.WriteString(fmt.Sprintf("  want: %#v\n", want))
	}
	sb.WriteString(fmt.Sprintf("   got: %#v\n", got))
	sb.WriteString(fmt.Sprintf("method: %s\n", method))
	sb.WriteString("\n")
	m.buildCallHistory(sb, "call history")
	sb.WriteString(hint)

	tb.Helper()
	m.fatal(index, sb.String())
}
