package meta

import (
	"context"
	"strings"
	"testing"
)

// ---

type target struct {
	td *targetTestDouble
}

type targetTestDouble struct {
	location string
	Full     *targetFull
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

	createdAtLocation := ""
	if m.td != nil && m.td.location != "" {
		createdAtLocation = m.td.location
	}
	panic(libMessageNotImplemented(
		"Target.Full", "Target.Full(ctx Context, id string) ([]Result, error)", "Full", libCallerLocation(2), createdAtLocation,
		"ctx", ctx, "input", input,
	))
}

type targetFull struct {
	Calls          []targetFullCall
	stub           func(ctx context.Context, input string) ([]Result, error)
	expects        []*targetFullExpect
	stubLocation   string
	verifyDisabled bool
}

func (m *targetFull) buildCallHistoryWithHeader(sb *strings.Builder) {
	if len(m.Calls) > 0 {
		sb.WriteString("call history:\n")
		m.buildCallHistory(sb)
	}
}

func (m *targetFull) buildCallHistory(sb *strings.Builder) {
	for i, call := range m.Calls {
		libMessageCallHistory(
			sb, i, m.expects[i].location, call.location,
			"ctx", call.Arguments.ctx, "input", call.Arguments.input,
		)
	}
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
	location := libCallerLocation(3)

	index := len(m.Calls)
	if index >= len(m.expects) {
		panic(libMessageTooManyCalls(
			"Target.Full", "Full", len(m.expects), index+1, location, m.buildCallHistory,
			"ctx", ctx, "input", input,
		))
	}

	expect := m.expects[index]
	if expect.matcher != nil {
		if !expect.matcher(ctx, input) {
			expect.tb.Helper()
			m.verifyDisabled = true
			expect.tb.Fatal(libMessageMatchFail(
				"Target.Full", expect.location, index+1, m.buildCallHistoryWithHeader,
				"ctx", ctx, "input", input,
			))
		}
	}

	if expect.arguments != nil {
		var msg string
		var ok bool

		ok, msg = libCompareByReflectEqual("ctx", "Target.Full", expect.location, index+1, expect.arguments.ctx, ctx, m.buildCallHistoryWithHeader)
		if !ok {
			expect.tb.Helper()
			m.verifyDisabled = true
			expect.tb.Fatal(msg)
		}

		ok, msg = libCompareByBasicComparison("input", "Target.Full", expect.location, index+1, expect.arguments.input, input, m.buildCallHistoryWithHeader)
		if !ok {
			expect.tb.Helper()
			m.verifyDisabled = true
			expect.tb.Fatal(msg)
		}
	}

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
	if matcher == nil {
		e.target.verifyDisabled = true
		e.target.expects[e.index].tb.Helper()
		e.target.expects[e.index].tb.Fatal(libMessageMatchByNil("Target.Full"))
	}

	e.target.expects[e.index].matcher = matcher
	e.target.expects[e.index].location = libCallerLocation(2)
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
	location := libCallerLocation(2)
	if stub == nil {
		panic(libMessageStubByNil("Target.Full", location))
	}

	if s.target.td.Full == nil {
		s.target.td.Full = &targetFull{}
	}

	if s.target.td.Full.stub != nil {
		s.target.td.Full.verifyDisabled = true
		panic(libMessageDuplicateStub("Target.Full", s.target.td.Full.stubLocation, location))
	}

	if len(s.target.td.Full.expects) > 0 {
		expect := s.target.td.Full.expects[0]
		s.target.td.Full.verifyDisabled = true
		panic(libMessageStubAfterExpect("Target.Full", expect.location, location))
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

	location := libCallerLocation(2)

	if e.target.td.Full.stub != nil {
		panic(libMessageExpectAfterStub("Target.Full", e.target.td.Full.stubLocation, location))
	}

	if tb == nil {
		panic(libMessageExpectByNil("Target.Full", "Full", location))
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
			tb.Fatal(libMessageExpectButNotCalled("Target.Full", len(expects), len(calls), index+1, e.target.td.Full.buildCallHistory))
		}
	})
	return &targetFullExpecter{index: index, target: e.target.td.Full}
}

func testTarget() *target {
	return &target{
		td: &targetTestDouble{
			location: libCallerLocation(2),
		},
	}
}
