package meta

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ---

type targetFull struct {
	Calls          []targetFullCall
	stub           func(ctx context.Context, input string) ([]Result, error)
	stubLocation   string
	verifyDisabled bool
	expects        []*targetFullExpect
}

func (m *targetFull) invokeStub(ctx context.Context, input string) ([]Result, error) {
	v0, v1 := m.stub(ctx, input)
	m.Calls = append(m.Calls, targetFullCall{
		Arguments: targetFullArgument{
			ctx:   ctx,
			input: input,
		},
		Returns: targetFullReturn{
			first:  v0,
			second: v1,
		},
	})
	return v0, v1
}

func (m *targetFull) invokeExpect(ctx context.Context, input string) ([]Result, error) {
	_, file, line, _ := runtime.Caller(2)                       // file = v0, line = v1
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line) // location = v2

	index := len(m.Calls) // index = v3
	if index >= len(m.expects) {
		sb := strings.Builder{}
		sb.WriteString("too many calls to Target.Full\n")
		sb.WriteString(fmt.Sprintf("\texpected: %d, got: %d\n\n", len(m.expects), index+1))
		for i, call := range m.Calls {
			sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", i+1, m.expects[i].location))
			sb.WriteString(fmt.Sprintf("\t   called at: %s\n", call.location))
			sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
			sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "ctx", call.Arguments.ctx))     // 5 is max of len(ctx), len(input)
			sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "input", call.Arguments.input)) // 5 is max of len(ctx), len(input)
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", index+1, "missing"))
		sb.WriteString(fmt.Sprintf("\t   called at: %s\n", location))
		sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
		sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "ctx", ctx))     // 5 is max of len(ctx), len(input)
		sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "input", input)) // 5 is max of len(ctx), len(input)
		sb.WriteString("\n")
		sb.WriteString("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().Full(t)\n")
		panic(sb.String())
	}

	expect := m.expects[index]
	if expect.matcher != nil {
		if !expect.matcher(ctx, input) {
			expect.tb.Helper()
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("Target.Full call #%d did not match\n", index+1))
			sb.WriteString(fmt.Sprintf("arguments:\n"))
			sb.WriteString(fmt.Sprintf("\t%5s = %#v\n", "ctx", ctx))     // 5 is max of len(ctx), len(input)
			sb.WriteString(fmt.Sprintf("\t%5s = %#v\n", "input", input)) // 5 is max of len(ctx), len(input)

			if len(m.Calls) > 0 {
				sb.WriteString("\ncall history:\n")
			}

			// duplicated - reduce later

			for i, call := range m.Calls {
				sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", i+1, m.expects[i].location))
				sb.WriteString(fmt.Sprintf("\t   called at: %s\n", call.location))
				sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
				sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "ctx", call.Arguments.ctx))     // 5 is max of len(ctx), len(input)
				sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "input", call.Arguments.input)) // 5 is max of len(ctx), len(input)
				sb.WriteString("\n")
			}

			sb.WriteString(fmt.Sprintf("hint: check the callback passed to Match at %s", expect.location))
			m.verifyDisabled = true
			expect.tb.Fatal(sb.String())
		}
	}
	//if m.expects[index].arguments != nil {
	//	// TODO: compare arguments
	//}

	m.Calls = append(m.Calls, targetFullCall{
		location: location,
		Arguments: targetFullArgument{
			ctx:   ctx,
			input: input,
		},
		Returns: targetFullReturn{
			first:  expect.returns.first,
			second: expect.returns.second,
		},
	})
	return expect.returns.first, expect.returns.second
}

type targetFullCall struct {
	location  string
	Arguments targetFullArgument
	Returns   targetFullReturn
}

type targetFullArgument struct {
	ctx   context.Context
	input string
}

type targetFullReturn struct {
	first  []Result
	second error
}

type targetFullExpect struct {
	matcher   func(ctx context.Context, input string) bool
	arguments *targetFullArgument
	returns   targetFullReturn
	location  string
	tb        testing.TB
}

type targetFullExpecter struct {
	index  int
	target *targetFull
}

func (e *targetFullExpecter) Return(first []Result, second error) {
	e.target.expects[e.index].returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

func (e *targetFullExpecter) Match(matcher func(ctx context.Context, input string) bool) *targetFullExpecterWithMatch {
	e.target.expects[e.index].matcher = matcher
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	e.target.expects[e.index].location = location
	return &targetFullExpecterWithMatch{index: e.index, target: e.target}
}

func (e *targetFullExpecter) CalledWith(ctx context.Context, input string) *targetFullExpecterWithArgs {
	e.target.expects[e.index].arguments = &targetFullArgument{
		ctx:   ctx,
		input: input,
	}
	return &targetFullExpecterWithArgs{index: e.index, target: e.target}
}

type targetFullExpecterWithArgs struct {
	index  int
	target *targetFull
}

func (e *targetFullExpecterWithArgs) Return(first []Result, second error) {
	e.target.expects[e.index].returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

type targetFullExpecterWithMatch struct {
	index  int
	target *targetFull
}

func (e *targetFullExpecterWithMatch) Return(first []Result, second error) {
	e.target.expects[e.index].returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

func (s *targetStubber) Full(stub func(ctx context.Context, input string) ([]Result, error)) *targetFull {
	if s.target.td.Full == nil {
		s.target.td.Full = &targetFull{}
	}
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	if s.target.td.Full.stub != nil {
		sb := strings.Builder{}
		sb.WriteString("duplicate STUB for Target.Full\n")
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "first used at", s.target.td.Full.stubLocation))
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "second used at", location))
		sb.WriteString("\thint: Target.Full is already stubbed, remove one of the above\n\n")
		s.target.td.Full.verifyDisabled = true
		panic(sb.String())
	}

	if len(s.target.td.Full.expects) > 0 {
		expect := s.target.td.Full.expects[0]
		sb := strings.Builder{}
		sb.WriteString("conflicting usage for Target.Full\n")
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "EXPECT used at", expect.location))
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "STUB used at", location))
		sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
		s.target.td.Full.verifyDisabled = true
		panic(sb.String())
	}

	s.target.td.Full.stub = stub
	s.target.td.Full.stubLocation = location
	return s.target.td.Full
}

func (e *targetExpecter) Full(tb testing.TB) *targetFullExpecter {
	if e.target.td == nil {
		e.target.td = &targetTestDouble{}
	}

	if e.target.td.Full == nil {
		e.target.td.Full = &targetFull{}
	}

	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	if e.target.td.Full.stub != nil {
		sb := strings.Builder{}
		sb.WriteString("conflicting usage for Target.Full\n")
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n", "STUB used at", e.target.td.Full.stubLocation))
		sb.WriteString(fmt.Sprintf("\t%14s: %s\n\n", "EXPECT used at", location))
		sb.WriteString("\thint: use either EXPECT or STUB for the same method, not both\n\n")
		panic(sb.String())
	}

	if tb == nil {
		sb := strings.Builder{}
		sb.WriteString("unexpected nil testing.TB in Target.Full\n")
		sb.WriteString(fmt.Sprintf("\tcalled at: %s\n\n", location))
		sb.WriteString("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n")
		sb.WriteString("\t\tspy := [var].STUB().Full(func(...) ...)\n")
		panic(sb.String())
	}

	e.target.td.Full.expects = append(e.target.td.Full.expects, &targetFullExpect{
		location: location,
		tb:       tb,
	})
	index := len(e.target.td.Full.expects) - 1

	tb.Helper()
	tb.Cleanup(func() {
		if e.target.td.Full.verifyDisabled {
			return
		}
		expects := e.target.td.Full.expects
		calls := e.target.td.Full.Calls
		if index >= len(calls) {
			tb.Helper()
			sb := strings.Builder{}
			sb.WriteString("Target.Full was not called as expected\n")
			sb.WriteString(fmt.Sprintf("\texpected: %d, got: %d\n\n", len(expects), len(calls)))
			for i, call := range calls {
				sb.WriteString(fmt.Sprintf("\t#%d expect at: %s\n", i+1, expects[i].location))
				sb.WriteString(fmt.Sprintf("\t   called at: %s\n", call.location))
				sb.WriteString(fmt.Sprintf("\t   arguments:\n"))
				sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "ctx", call.Arguments.ctx))     // 5 is max of len(ctx), len(input)
				sb.WriteString(fmt.Sprintf("\t\t%5s = %#v\n", "input", call.Arguments.input)) // 5 is max of len(ctx), len(input)
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("\t#%d never called\n\n", index+1))
			sb.WriteString("\thint: add the missing call or remove the EXPECT above")
			tb.Fatal(sb.String())
		}
	})
	return &targetFullExpecter{index: index, target: e.target.td.Full}
}

// ---

func testTarget() *target {
	_, file, line, _ := runtime.Caller(1)
	return &target{
		td: &targetTestDouble{
			location: fmt.Sprintf("%s:%d", filepath.Base(file), line),
		},
	}
}
