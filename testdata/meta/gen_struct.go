package meta

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type target struct {
	td *targetTestDouble
}

type targetStubber struct {
	target *target
}

func (m *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

type targetExpecter struct {
	target *target
}

func (m *target) EXPECT() *targetExpecter {
	return &targetExpecter{target: m}
}

type targetTestDouble struct {
	location string
	Full     *targetFull
}

func (m *target) Full(ctx context.Context, input string) ([]Result, error) {
	if m.td != nil && m.td.Full != nil && m.td.Full.stub != nil {
		return m.td.Full.stub(ctx, input)
	}

	if m.td != nil && m.td.Full != nil && len(m.td.Full.expects) > 0 {
		index := len(m.td.Full.Calls)
		if index < len(m.td.Full.expects) {
			m.td.Full.expects[index].tb.Helper()
		}
		return m.td.Full.invokeExpect(ctx, input)
	}

	sb := strings.Builder{}
	sb.WriteString("unexpected call to Target.Full\n")
	sb.WriteString("signature: Target.Full(ctx Context, id string) ([]Result, error)\n")
	_, file, line, _ := runtime.Caller(1)
	calledAt := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	sb.WriteString(fmt.Sprintf("called at: %s\n", calledAt))

	sb.WriteString("arguments:\n")
	sb.WriteString(fmt.Sprintf("\t%5s = %#v\n", "ctx", ctx))     // 5 is max of len(ctx), len(input)
	sb.WriteString(fmt.Sprintf("\t%5s = %#v\n", "input", input)) // 5 is max of len(ctx), len(input)

	location := ""
	if m.td != nil && m.td.location != "" {
		location = " after " + m.td.location
	}
	sb.WriteString(fmt.Sprintf("\nhint:%s use one of:\n\t[var].EXPECT().Full(t)\n\t[var].STUB().Full(func(...) ...)\n\n", location))
	panic(sb.String())
}
