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
	Calls        []targetFullCall
	stub         func(ctx context.Context, input string) ([]Result, error)
	stubLocation string
	expects      []*targetFullExpect
	verified     bool
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
			m.verified = true
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
			m.verified = true
			expect.tb.Fatal(msg)
		}

		ok, msg = libCompareByBasicComparison("input", "Target.Full", expect.location, index+1, expect.arguments.input, input, m.buildCallHistoryWithHeader)
		if !ok {
			expect.tb.Helper()
			m.verified = true
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

func (m *targetFull) verify(index int) {
	if m.verified {
		return
	}

	if index >= len(m.Calls) {
		m.expects[index].tb.Helper()
		m.expects[index].tb.Fatal(libMessageExpectButNotCalled("Target.Full", len(m.expects), len(m.Calls), index+1, m.buildCallHistory))
	}
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
	target *targetFull
	expect *targetFullExpect
}

func (e *targetFullExpecter) Return(first []Result, second error) {
	e.expect.returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

func (e *targetFullExpecter) Match(matcher func(ctx context.Context, input string) bool) *targetFullExpecterWithMatch {
	if matcher == nil {
		e.target.verified = true
		e.expect.tb.Helper()
		e.expect.tb.Fatal(libMessageMatchByNil("Target.Full"))
	}

	e.expect.matcher = matcher
	e.expect.location = libCallerLocation(2)
	return &targetFullExpecterWithMatch{expect: e.expect}
}

func (e *targetFullExpecter) CalledWith(ctx context.Context, input string) *targetFullExpecterWithArgs {
	e.expect.arguments = &targetFullArgument{
		ctx:   ctx,
		input: input,
	}
	return &targetFullExpecterWithArgs{expect: e.expect}
}

type targetFullExpecterWithArgs struct {
	expect *targetFullExpect
}

func (e *targetFullExpecterWithArgs) Return(first []Result, second error) {
	e.expect.returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

type targetFullExpecterWithMatch struct {
	expect *targetFullExpect
}

func (e *targetFullExpecterWithMatch) Return(first []Result, second error) {
	e.expect.returns = targetFullReturn{
		first:  first,
		second: second,
	}
}

func (s *targetStubber) Full(stub func(ctx context.Context, input string) ([]Result, error)) *targetFull {
	location := libCallerLocation(2)
	if stub == nil {
		panic(libMessageStubByNil("Target.Full", location))
	}

	var spy = s.target.td.Full
	if spy == nil {
		spy = &targetFull{stubLocation: location}
		s.target.td.Full = spy
	}

	if spy.stub != nil {
		spy.verified = true
		panic(libMessageDuplicateStub("Target.Full", spy.stubLocation, location))
	}

	if len(spy.expects) > 0 {
		spy.verified = true
		panic(libMessageStubAfterExpect("Target.Full", spy.expects[0].location, location))
	}

	spy.stub = stub
	return spy
}

func (e *targetExpecter) Full(tb testing.TB) *targetFullExpecter {
	if e.target.td == nil {
		e.target.td = &targetTestDouble{}
	}

	var mock = e.target.td.Full
	if mock == nil {
		mock = &targetFull{}
		e.target.td.Full = mock
	}

	location := libCallerLocation(2)
	if mock.stub != nil {
		panic(libMessageExpectAfterStub("Target.Full", mock.stubLocation, location))
	}
	if tb == nil {
		panic(libMessageExpectByNil("Target.Full", "Full", location))
	}

	mock.expects = append(mock.expects, &targetFullExpect{
		location: location,
		tb:       tb,
	})

	index := len(mock.expects) - 1
	tb.Helper()
	tb.Cleanup(func() {
		tb.Helper()
		mock.verify(index)
	})

	return &targetFullExpecter{target: mock, expect: mock.expects[index]}
}

func testTarget() *target {
	return &target{
		td: &targetTestDouble{
			location: libCallerLocation(2),
		},
	}
}
