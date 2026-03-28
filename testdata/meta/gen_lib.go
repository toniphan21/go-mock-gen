package meta

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func targetCallerLocation(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func targetMessageWriteArguments(sb *strings.Builder, template string, args []any) {
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

func targetMessageNotImplemented(target, signature, method, calledAt, createdLocation string, args ...any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected call to %s\n", target))
	sb.WriteString(fmt.Sprintf("signature: %s\n", signature))
	sb.WriteString(fmt.Sprintf("called at: %s\n", calledAt))

	sb.WriteString("arguments:\n")
	targetMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)

	location := ""
	if createdLocation != "" {
		location = " after " + createdLocation
	}
	sb.WriteString(fmt.Sprintf(
		"\nhint:%s use one of:\n\t[var].EXPECT().%s(t)\n\t[var].STUB().%s(func(...) ...)\n\n",
		location, method, method,
	))
	return sb.String()
}

func targetMessageCallHistory(sb *strings.Builder, index int, expectedAt, calledAt string, args ...any) string {
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", index+1, expectedAt))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", calledAt))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	targetMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	return sb.String()
}

func targetMessageTooManyCalls(target, method string, want, got int, calledAt string, fn func(sb *strings.Builder), args ...any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("too many calls to %s\n", target))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	fn(sb)
	sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", got, "missing"))
	sb.WriteString(fmt.Sprintf("\t   called at: %s\n", calledAt))
	sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
	targetMessageWriteArguments(sb, "\t\t%[MAX-KEY-LEN]s = %#v\n", args)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().%s(t)\n", method))
	return sb.String()
}

func targetMessageMatchByNil(target string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s Match received a nil function\n", target))
	sb.WriteString("\thint: provide a valid function")
	return sb.String()
}

func targetMessageStubByNil(target, calledAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s STUB received a nil function\n", target))
	sb.WriteString(fmt.Sprintf("called at: %s\n\n", calledAt))
	sb.WriteString("hint: provide a valid function\n")
	return sb.String()
}

func targetMessageExpectByNil(target, method, calledAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected nil testing.TB in %s\n", target))
	sb.WriteString(fmt.Sprintf("\tcalled at: %s\n\n", calledAt))
	sb.WriteString("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n")
	sb.WriteString(fmt.Sprintf("\t\tspy := [var].STUB().%s(func(...) ...)\n", method))
	panic(sb.String())
}

func targetMessageExpectAfterStub(target, stubAt, expectAt string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("conflicting usage for %s\n", target))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "STUB used at", stubAt))
	sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "EXPECT used at", expectAt))
	sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
	return sb.String()
}

func targetMessageExpectButNotCalled(target string, want, got, expectNo int, fn func(sb *strings.Builder)) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s was not called as expected\n", target))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	fn(&sb)
	sb.WriteString(fmt.Sprintf("\t#%d never called\n\n", expectNo))
	sb.WriteString("\thint: add the missing call or remove the EXPECT above")
	return sb.String()
}
