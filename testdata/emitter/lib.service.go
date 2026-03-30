package emitter

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func serviceCallerLocation(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

type serviceMockMethod interface {
	methodName() string
	interfaceName() string
	buildCallHistory(sb *strings.Builder, header string)
	fatal(index int, msg string)
	panic(msg string)
}

func serviceMessageWriteArguments(sb *strings.Builder, template string, args []any) {
	maxLen := 0
	for i := 0; i < len(args); i += 2 {
		str, ok := args[i].(string)
		if !ok {
			str = fmt.Sprintf("%v", args[i])
		}
		maxLen = max(maxLen, len(str))
	}

	format := strings.ReplaceAll(template, "[MAX-KEY-LEN]", strconv.Itoa(maxLen))
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", args[i])
		}

		var val any
		if i+1 < len(args) {
			val = args[i+1]
		}
		sb.WriteString(fmt.Sprintf(format, key, val))
	}
}

func serviceMessageMatchFail(m serviceMockMethod, matchedAt string, index int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s call #%d did not match\n", m.interfaceName(), m.methodName(), index+1))
	sb.WriteString(fmt.Sprintf("arguments:\n"))
	serviceMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	m.buildCallHistory(sb, "call history")
	sb.WriteString(fmt.Sprintf("hint: check the callback passed to Match at %s", matchedAt))
	return sb.String()
}

func serviceMessageNotImplemented(interfaceName string, methodName string, signature string, createdLocation string, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected call to %s.%s\n", interfaceName, methodName))
	sb.WriteString(fmt.Sprintf("signature: %s.%s%s\n", interfaceName, methodName, signature))
	sb.WriteString(fmt.Sprintf("called at: %s\n", serviceCallerLocation(3)))
	sb.WriteString("arguments:\n")
	serviceMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)

	location := ""
	if createdLocation != "" {
		location = " after " + createdLocation
	}
	sb.WriteString(fmt.Sprintf("\nhint:%s use one of:\n\t[var].EXPECT().%s(t)\n\t[var].STUB().%s(func(...) ...)\n\n", location, methodName, methodName))
	return sb.String()
}

func serviceMessageCallHistory(sb *strings.Builder, index int, expectedAt string, calledAt string, args []any) string {
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", index+1, expectedAt))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", calledAt))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	serviceMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	return sb.String()
}

func serviceMessageTooManyCalls(m serviceMockMethod, want int, got int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("too many calls to %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", got, "missing"))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", serviceCallerLocation(4)))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	serviceMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().%s(t)\n", m.methodName()))
	return sb.String()
}

func serviceMessageMatchByNil(m serviceMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s Match received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}

func serviceMessageExpectByNil(m serviceMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected nil testing.TB in %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\tcalled at: %s\n\n", serviceCallerLocation(3)))
	sb.WriteString("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n")
	sb.WriteString(fmt.Sprintf("\t\tspy := [var].STUB().%s(func(...) ...)\n", m.methodName()))
	return sb.String()
}

func serviceMessageExpectAfterStub(m serviceMockMethod, stubAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "STUB used at", stubAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "EXPECT used at", serviceCallerLocation(3)))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}

func serviceMessageStubByNil(m serviceMockMethod, calledAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s STUB received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("called at: %s\n\n", calledAt))
	sb.WriteString("hint: provide a valid function\n")
	return sb.String()
}

func serviceMessageStubAfterExpect(m serviceMockMethod, expectAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "EXPECT used at", expectAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "STUB used at", serviceCallerLocation(3)))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}

func serviceMessageDuplicateStub(m serviceMockMethod, firstUsedAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("duplicate STUB for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "first used at", firstUsedAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "second used at", serviceCallerLocation(3)))
	sb.WriteString(fmt.Sprintf("\thint: %s.%s is already stubbed, remove one of the above\n\n", m.interfaceName(), m.methodName()))
	return sb.String()
}

func serviceMessageExpectButNotCalled(m serviceMockMethod, want int, got int, index int) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s was not called as expected\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d never called\n\n", index+1))
	sb.WriteString("\thint: add the missing call or remove the EXPECT above")
	return sb.String()
}

func serviceMessageMatchArgByNil(m serviceMockMethod, method string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s %s received a nil function\n", m.interfaceName(), m.methodName(), method))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}

func serviceMessageDuplicateMatchArg(m serviceMockMethod, method string, firstUsedAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("duplicate %s for %s.%s\n", method, m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "first used at", firstUsedAt))
	sb.WriteString("\thint: each argument can only be matched once, remove one of the above")
	return sb.String()
}

func serviceMessageMatchArgHint() string {
	return fmt.Sprintf("\thint: check argument matching at %s\n\t\tor use STUB for fine-grained control", serviceCallerLocation(3))
}

func serviceMatchArgument[T any](m serviceMockMethod, index int, name string, got T, match func(T) bool, wants map[string]any, methods map[string]string, hints map[string]string, tb testing.TB, expectAt string) {
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

func serviceReflectEqualMatcher[T any](want T) func(T) bool {
	return func(got T) bool {
		return reflect.DeepEqual(want, got)
	}
}

func serviceBasicComparisonMatcher[T comparable](want T) func(T) bool {
	return func(got T) bool {
		return want == got
	}
}
