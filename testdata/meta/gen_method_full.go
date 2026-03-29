package meta

import (
	"context"
	"fmt"
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
	interfaceName, methodName, signature := "Target", "Full", "(ctx Context, id string) ([]Result, error)"
	args := []any{"ctx", ctx, "input", input}

	if m.td == nil {
		panic(libMessageNotImplemented(interfaceName, methodName, signature, "", args))
	}

	if mock := m.td.Full; mock != nil {
		switch {
		case mock.stub != nil:
			return mock.invokeStub(ctx, input)

		case len(mock.expects) > 0:
			index := len(mock.Calls)
			if index < len(mock.expects) {
				mock.expects[index].tb.Helper()
			}
			return mock.invokeExpect(ctx, input)
		}
	}
	panic(libMessageNotImplemented(interfaceName, methodName, signature, m.td.location, args))
}

type targetFull struct {
	Calls        []targetFullCall
	stub         func(ctx context.Context, input string) ([]Result, error)
	stubLocation string
	expects      []*targetFullExpect
	verified     bool
}

func (m *targetFull) methodName() string {
	return "Full"
}

func (m *targetFull) interfaceName() string {
	return "Target"
}

func (m *targetFull) fatal(index int, msg string) {
	m.verified = true
	m.expects[index].tb.Helper()
	m.expects[index].tb.Fatal(msg)
}

func (m *targetFull) panic(msg string) {
	m.verified = true
	panic(msg)
}

func (m *targetFull) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(m.Calls) != 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, call := range m.Calls {
		args := []any{"ctx", call.Arguments.ctx, "input", call.Arguments.input}
		libMessageCallHistory(sb, i, m.expects[i].location, call.location, args)
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
	args := []any{"ctx", ctx, "input", input}
	index := len(m.Calls)
	if index >= len(m.expects) {
		panic(libMessageTooManyCalls(m, len(m.expects), index+1, args))
	}

	expect := m.expects[index]
	if expect.matcher != nil && !expect.matcher(ctx, input) {
		expect.tb.Helper()
		m.fatal(index, libMessageMatchFail(m, expect.location, index, args))
	}

	if expect.arguments != nil {
		expect.tb.Helper()
		libCompareByReflectEqual(m, "ctx", expect.arguments.ctx, ctx, expect.tb, expect.location, index)
		libCompareByBasicComparison(m, "input", expect.arguments.input, input, expect.tb, expect.location, index)
	}

	m.Calls = append(m.Calls, targetFullCall{
		location: libCallerLocation(3),
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
	index     int
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
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchByNil(e.target))
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
	if s.target.td == nil {
		s.target.td = &targetTestDouble{}
	}

	var spy = s.target.td.Full
	if spy == nil {
		spy = &targetFull{stubLocation: libCallerLocation(2)}
		s.target.td.Full = spy
	}

	if stub == nil {
		spy.panic(libMessageStubByNil(spy, libCallerLocation(2)))
	}

	if spy.stub != nil {
		spy.panic(libMessageDuplicateStub(spy, spy.stubLocation))
	}

	if len(spy.expects) > 0 {
		spy.panic(libMessageStubAfterExpect(spy, spy.expects[0].location))
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

	if mock.stub != nil {
		panic(libMessageExpectAfterStub(mock, mock.stubLocation))
	}

	if tb == nil {
		panic(libMessageExpectByNil(mock))
	}

	index := len(mock.expects)
	mock.expects = append(mock.expects, &targetFullExpect{
		location: libCallerLocation(2),
		index:    index,
		tb:       tb,
	})

	tb.Helper()
	tb.Cleanup(func() {
		if !mock.verified && index >= len(mock.Calls) {
			mock.expects[index].tb.Helper()
			mock.fatal(index, libMessageExpectButNotCalled(mock, len(mock.expects), len(mock.Calls), index))
		}
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
