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
	 * sb.WriteString(fmt.Sprintf("\t%[MAX-KEY-LEN]s = %#v\n", args[0,2,4...], args[1,3,5...]))
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

func targetMessageNotImplemented(target, signature, calledAt, createdLocation string, args ...any) string {
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
	sb.WriteString(fmt.Sprintf("\nhint:%s use one of:\n\t[var].EXPECT().Full(t)\n\t[var].STUB().Full(func(...) ...)\n\n", location))
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
