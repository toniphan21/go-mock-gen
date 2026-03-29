package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageNotImplemented(interfaceName string, methodName string, signature string, createdLocation string, args []any) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("unexpected call to %s.%s\n", interfaceName, methodName))
	sb.WriteString(fmt.Sprintf("signature: %s.%s%s\n", interfaceName, methodName, signature))
	sb.WriteString(fmt.Sprintf("called at: %s\n", repositoryCallerLocation(3)))
	sb.WriteString("arguments:\n")
	repositoryMessageWriteArguments(sb, "\t%[MAX-KEY-LEN]s = %#v\n", args)

	location := ""
	if createdLocation != "" {
		location = " after " + createdLocation
	}
	sb.WriteString(fmt.Sprintf("\nhint:%s use one of:\n\t[var].EXPECT().%s(t)\n\t[var].STUB().%s(func(...) ...)\n\n", location, methodName, methodName))
	return sb.String()
}
