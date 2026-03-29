package meta

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type libMockMethod interface {
	methodName() string
	interfaceName() string
	buildCallHistory(sb *strings.Builder, header string)
	fatal(index int, msg string)
	panic(msg string)
}

func libMessageMatchFail(m libMockMethod, matchedAt string, index int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s call #%d did not match\n", m.interfaceName(), m.methodName(), index+1))
	sb.WriteString(fmt.Sprintf("arguments:\n"))
	libMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	m.buildCallHistory(sb, "call history")
	sb.WriteString(fmt.Sprintf("hint: check the callback passed to Match at %s", matchedAt))
	return sb.String()
}

func libCompareByReflectEqual[M libMockMethod, T any](m M, argName string, want T, got T, tb testing.TB, expectAt string, index int) {
	if reflect.DeepEqual(want, got) {
		return
	}

	tb.Helper()
	msg := libMessageArgumentMismatched(
		argName,
		m.interfaceName()+"."+m.methodName(),
		expectAt,
		"reflect.DeepEqual",
		index+1,
		want, got,
		func(sb *strings.Builder) {
			m.buildCallHistory(sb, "call history")
		},
	)
	m.fatal(index, msg)
}

func libCompareByBasicComparison[M libMockMethod, T comparable](m M, argName string, want T, got T, tb testing.TB, expectAt string, index int) {
	if want == got {
		return
	}

	tb.Helper()
	msg := libMessageArgumentMismatched(
		argName,
		m.interfaceName()+"."+m.methodName(),
		expectAt,
		"==",
		index+1,
		want,
		got,
		func(sb *strings.Builder) {
			m.buildCallHistory(sb, "call history")
		},
	)
	m.fatal(index, msg)
}

func libCallerLocation(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func libMessageWriteArguments(sb *strings.Builder, template string, args []any) {
	/*
	 * MAX-KEY-LEN is max of len(args[0,2,4,6...))
	 * sb.WriteString(fmt.Sprintf(strings.ReplaceAll(template, "[MAX-KEY-LEN]"), args[0,2,4...], args[1,3,5...]))
	 */
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

func libMessageNotImplemented(interfaceName, methodName, signature, createdLocation string, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected call to %s.%s\n", interfaceName, methodName))
	sb.WriteString(fmt.Sprintf("signature: %s.%s%s\n", interfaceName, methodName, signature))
	sb.WriteString(fmt.Sprintf("called at: %s\n", libCallerLocation(3)))

	sb.WriteString("arguments:\n")
	libMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)

	location := ""
	if createdLocation != "" {
		location = " after " + createdLocation
	}
	sb.WriteString(fmt.Sprintf(
		"\nhint:%s use one of:\n\t[var].EXPECT().%s(t)\n\t[var].STUB().%s(func(...) ...)\n\n",
		location, methodName, methodName,
	))
	return sb.String()
}

func libMessageCallHistory(sb *strings.Builder, index int, expectedAt, calledAt string, args []any) string {
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", index+1, expectedAt))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", calledAt))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	libMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	return sb.String()
}

func libMessageTooManyCalls(m libMockMethod, want, got int, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("too many calls to %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", got, "missing"))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", libCallerLocation(4)))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	libMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().%s(t)\n", m.methodName()))
	return sb.String()
}

func libMessageMatchByNil(m libMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s Match received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}

func libMessageArgumentMismatched(argName, target, expectAt string, comparedBy string, callNo int, want any, got any, fn func(builder *strings.Builder)) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s call #%d argument \"%s\" did not match\n", target, callNo, argName))
	sb.WriteString(fmt.Sprintf("  want: %#v\n", want))
	sb.WriteString(fmt.Sprintf("   got: %#v\n", got))
	sb.WriteString(fmt.Sprintf("method: %s\n", comparedBy))
	sb.WriteString("\n")
	fn(sb)
	sb.WriteString(fmt.Sprintf("hint: for custom matching use .Match(func(...) bool) at %s\n\tor use STUB for fine-grained control", expectAt))
	return sb.String()
}

func libMessageExpectByNil(m libMockMethod) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected nil testing.TB in %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\tcalled at: %s\n\n", libCallerLocation(3)))
	sb.WriteString("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n")
	sb.WriteString(fmt.Sprintf("\t\tspy := [var].STUB().%s(func(...) ...)\n", m.methodName()))
	panic(sb.String())
}

func libMessageExpectAfterStub(m libMockMethod, stubAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "STUB used at", stubAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "EXPECT used at", libCallerLocation(3)))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}

func libMessageStubByNil(m libMockMethod, calledAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s STUB received a nil function\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("called at: %s\n\n", calledAt))
	sb.WriteString("hint: provide a valid function\n")
	return sb.String()
}

func libMessageStubAfterExpect(m libMockMethod, expectAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "EXPECT used at", expectAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "STUB used at", libCallerLocation(3)))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}

func libMessageDuplicateStub(m libMockMethod, firstUsedAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("duplicate STUB for %s.%s\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "first used at", firstUsedAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "second used at", libCallerLocation(3)))
	sb.WriteString(fmt.Sprintf("\thint: %s.%s is already stubbed, remove one of the above\n\n", m.interfaceName(), m.methodName()))
	return sb.String()
}

func libMessageExpectButNotCalled(m libMockMethod, want, got, index int) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s was not called as expected\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d never called\n\n", index+1))
	sb.WriteString("\thint: add the missing call or remove the EXPECT above")
	return sb.String()
}
