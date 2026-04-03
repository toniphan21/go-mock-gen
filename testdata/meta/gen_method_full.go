package meta

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// ---

// done - TargetData.constructorCode()
func testTarget() *target {
	return &target{td: &targetTestDouble{location: libCallerLocation(2)}}
}

// done - TargetData.targetStructCode()
type target struct {
	td *targetTestDouble
}

// done - TargetData.targetBuiltinFuncCode()
func (m *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

// done - TargetData.GenerateCode()
func (m *target) EXPECT() *targetExpecter { // skip:!expect
	return &targetExpecter{target: m}
}

// done - TargetData.implementationCode()
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

		case len(mock.expects) > 0: // skip:!expect
			index := len(mock.Calls)
			if index < len(mock.expects) {
				mock.expects[index].tb.Helper()
			}
			return mock.invokeExpect(ctx, input)
		}
	}
	panic(libMessageNotImplemented(interfaceName, methodName, signature, m.td.location, args))
}

// done - TargetData.targetStructCode()
type targetTestDouble struct {
	location string
	Full     *targetFull
}

// done - TargetStubberData.targetStubberStructCode()
type targetStubber struct {
	target *target
}

// done - TargetStubberData.stubCode()
func (s *targetStubber) Full(stub func(ctx context.Context, input string) ([]Result, error)) *targetFull {
	if s.target.td == nil {
		s.target.td = &targetTestDouble{}
	}

	spy := s.target.td.Full
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

	if len(spy.expects) > 0 { // skip:!expect
		spy.panic(libMessageStubAfterExpect(spy, spy.expects[0].location))
	}

	spy.stub = stub
	return spy
}

// done - TargetExpecterData.targetExpecterStructCode()
type targetExpecter struct { // skip:!expect
	target *target
}

// done - TargetExpecterData.expectCode()
func (e *targetExpecter) Full(tb testing.TB) *targetFullExpecter { // skip:!expect
	if e.target.td == nil {
		e.target.td = &targetTestDouble{}
	}

	var mock = e.target.td.Full
	if mock == nil {
		mock = &targetFull{}
		e.target.td.Full = mock
	}

	if mock.stub != nil {
		mock.panic(libMessageExpectAfterStub(mock, mock.stubLocation))
	}

	if tb == nil {
		mock.panic(libMessageExpectByNil(mock))
	}

	index := len(mock.expects)
	mock.expects = append(mock.expects, &targetFullExpect{
		location:         libCallerLocation(2),
		matcher:          &targetFullArgumentMatcher{},
		matcherWants:     make(map[string]any),
		matcherMethods:   make(map[string]string),
		matcherHints:     make(map[string]string),
		matcherLocations: make(map[string]string),
		index:            index,
		tb:               tb,
	})

	tb.Helper()
	tb.Cleanup(func() { tb.Helper(); mock.verify(index) })

	return &targetFullExpecter{target: mock, expect: mock.expects[index]}
}

// done - MethodData.structCode()
type targetFull struct {
	Calls        []targetFullCall
	stub         func(ctx context.Context, input string) ([]Result, error)
	stubLocation string
	expects      []*targetFullExpect // skip:!expect
	verified     bool                // skip:!expect
}

// done - MethodData.methodNameFuncCode()
func (m *targetFull) methodName() string {
	return "Full"
}

// done - MethodData.interfaceNameFuncCode()
func (m *targetFull) interfaceName() string {
	return "Target"
}

// done - MethodData.fatalFuncCode()
func (m *targetFull) fatal(index int, msg string) {
	m.verified = true              // skip:!expect
	m.expects[index].tb.Helper()   // skip:!expect
	m.expects[index].tb.Fatal(msg) // skip:!expect
}

// done - MethodData.panicFuncCode()
func (m *targetFull) panic(msg string) {
	m.verified = true // skip:!expect
	panic(msg)
}

// done - MethodData.buildCallHistoryFuncCode()
func (m *targetFull) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(m.Calls) != 0 { // skip:!expect
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, call := range m.Calls { // skip:!expect
		args := []any{"ctx", call.Argument.ctx, "input", call.Argument.input}
		libMessageCallHistory(sb, i, m.expects[i].location, call.Location, args)
	}
}

// done - MethodData.invokeStubFuncCode()
func (m *targetFull) invokeStub(ctx context.Context, input string) ([]Result, error) {
	v0, v1 := m.stub(ctx, input)
	return m.capture(
		targetFullArgument{ctx: ctx, input: input},
		targetFullReturn{first: v0, second: v1},
	)
}

// done - MethodData.invokeExpectFuncCode()
func (m *targetFull) invokeExpect(ctx context.Context, input string) ([]Result, error) { // skip:!expect
	args := []any{"ctx", ctx, "input", input}
	index := len(m.Calls)
	if index >= len(m.expects) {
		m.panic(libMessageTooManyCalls(m, len(m.expects), index+1, args))
	}

	expect := m.expects[index]
	if expect.match != nil && !expect.match(ctx, input) {
		expect.tb.Helper()
		m.fatal(index, libMessageMatchFail(m, expect.matchLocation, index, args))
	}

	expect.tb.Helper()
	libMatchArgument(m, index, "ctx", ctx, expect.matcher.ctx, expect.matcherWants, expect.matcherMethods, expect.matcherHints, expect.tb, expect.matcherLocations["ctx"])
	libMatchArgument(m, index, "input", input, expect.matcher.input, expect.matcherWants, expect.matcherMethods, expect.matcherHints, expect.tb, expect.matcherLocations["input"])

	return m.capture(
		targetFullArgument{ctx: ctx, input: input},
		expect.returns,
	)
}

// done - MethodData.captureFuncCode()
func (m *targetFull) capture(args targetFullArgument, returns targetFullReturn) ([]Result, error) {
	m.Calls = append(m.Calls, targetFullCall{
		Location: libCallerLocation(4),
		Argument: args,
		Return:   returns,
	})
	return returns.first, returns.second
}

// done - MethodData.verifyFuncCode()
func (m *targetFull) verify(index int) { // skip:!expect
	if !m.verified && index >= len(m.Calls) {
		m.expects[index].tb.Helper()
		m.expects[index].tb.Fatal(libMessageExpectButNotCalled(m, len(m.expects), len(m.Calls), index))
	}
}

// done - MethodData.callStructCode()
type targetFullCall struct {
	Location string
	Argument targetFullArgument
	Return   targetFullReturn
}

// done - MethodData.argumentStructCode()
type targetFullArgument struct {
	ctx   context.Context
	input string
}

// done - MethodData.argumentMatcherStructCode()
type targetFullArgumentMatcher struct { // skip:!expect
	ctx   func(context.Context) bool
	input func(string) bool
}

// done - MethodData.returnStructCode()
type targetFullReturn struct {
	first  []Result
	second error
}

// done - MethodData.expectStructCode()
type targetFullExpect struct { // skip:!expect
	match            func(ctx context.Context, input string) bool
	matchLocation    string
	matcher          *targetFullArgumentMatcher
	matcherWants     map[string]any
	matcherMethods   map[string]string
	matcherHints     map[string]string
	matcherLocations map[string]string
	returns          targetFullReturn
	location         string
	index            int
	tb               testing.TB
}

// done - MethodExpecterData.structCode()
type targetFullExpecter struct { // skip:!expect
	target *targetFull
	expect *targetFullExpect
}

// done - MethodExpecterData.GenerateCode()
func (e *targetFullExpecter) Return(first []Result, second error) { // skip:!expect
	e.expect.returns = targetFullReturn{first: first, second: second}
}

// done - MethodExpecterData.matchFuncCode()
func (e *targetFullExpecter) Match(matcher func(ctx context.Context, input string) bool) *targetFullExpecterWithMatch { // skip:!expect
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchByNil(e.target))
	}

	e.expect.match = matcher
	e.expect.matchLocation = libCallerLocation(2)
	return &targetFullExpecterWithMatch{expect: e.expect}
}

// done - MethodExpecterData.argumentFuncCode()
func (e *targetFullExpecter) MatchCtx(matcher func(ctx context.Context) bool) *targetFullExpecterWithMatchArg { // skip:!expect
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchCtx"))
	}

	e.expect.matcher.ctx = matcher
	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
	e.expect.matcherHints["ctx"] = libMessageMatchArgHint()
	return &targetFullExpecterWithMatchArg{expect: e.expect, target: e.target}
}

// done - MethodExpecterData.argumentFuncCode()
func (e *targetFullExpecter) MatchInput(matcher func(input string) bool) *targetFullExpecterWithMatchArg { // skip:!expect
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchInput"))
	}

	e.expect.matcher.input = matcher
	e.expect.matcherLocations["input"] = libCallerLocation(2)
	e.expect.matcherHints["input"] = libMessageMatchArgHint()
	return &targetFullExpecterWithMatchArg{expect: e.expect, target: e.target}
}

// done - MethodExpecterData.withFuncCode()
func (e *targetFullExpecter) With(ctx context.Context, input string) *targetFullExpecterWithValue { // skip:!expect
	e.WithCtx(ctx)
	e.expect.matcherLocations["ctx"] = libCallerLocation(2)

	e.WithInput(input)
	e.expect.matcherLocations["input"] = libCallerLocation(2)

	return &targetFullExpecterWithValue{expect: e.expect}
}

// done - MethodExpecterData.argumentFuncCode()
func (e *targetFullExpecter) WithCtx(ctx context.Context) *targetFullExpecterWithValueArg { // skip:!expect
	e.expect.matcher.ctx = libReflectEqualMatcher(ctx)
	e.expect.matcherWants["ctx"] = ctx
	e.expect.matcherMethods["ctx"] = "reflect.DeepEqual"
	e.expect.matcherLocations["ctx"] = libCallerLocation(2)

	return &targetFullExpecterWithValueArg{expect: e.expect, target: e.target}
}

// done - MethodExpecterData.argumentFuncCode()
func (e *targetFullExpecter) WithInput(input string) *targetFullExpecterWithValueArg { // skip:!expect
	e.expect.matcher.input = libBasicComparisonMatcher(input)
	e.expect.matcherWants["input"] = input
	e.expect.matcherMethods["input"] = "=="
	e.expect.matcherLocations["input"] = libCallerLocation(2)

	return &targetFullExpecterWithValueArg{expect: e.expect, target: e.target}
}

// done - MethodExpecterValueData.structCode()
type targetFullExpecterWithValue struct { // skip:!expect
	expect *targetFullExpect
}

// done - MethodExpecterValueData.GenerateCode()
func (e *targetFullExpecterWithValue) Return(first []Result, second error) { // skip:!expect
	e.expect.returns = targetFullReturn{first: first, second: second}
}

// done - MethodExpecterMatchData.structCode()
type targetFullExpecterWithMatch struct { // skip:!expect
	expect *targetFullExpect
}

// done - MethodExpecterMatchData.GenerateCode()
func (e *targetFullExpecterWithMatch) Return(first []Result, second error) { // skip:!expect
	e.expect.returns = targetFullReturn{first: first, second: second}
}

// done - MethodExpecterMatchArgData.structCode()
type targetFullExpecterWithMatchArg struct { // skip:!expect
	expect *targetFullExpect
	target *targetFull
}

// done - MethodExpecterMatchArgData.GenerateCode()
func (e *targetFullExpecterWithMatchArg) Return(first []Result, second error) { // skip:!expect
	e.expect.returns = targetFullReturn{first: first, second: second}
}

// done - MethodExpecterMatchArgData.GenerateCode()
func (e *targetFullExpecterWithMatchArg) MatchCtx(matcher func(ctx context.Context) bool) *targetFullExpecterWithMatchArg { // skip:!expect
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchCtx"))
	}

	if e.expect.matcher.ctx != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchCtx", e.expect.matcherLocations["ctx"]))
	}

	e.expect.matcher.ctx = matcher
	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
	e.expect.matcherHints["ctx"] = libMessageMatchArgHint()
	return e
}

// done - MethodExpecterMatchArgData.GenerateCode()
func (e *targetFullExpecterWithMatchArg) MatchInput(matcher func(input string) bool) *targetFullExpecterWithMatchArg { // skip:!expect
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchInput"))
	}

	if e.expect.matcher.input != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchInput", e.expect.matcherLocations["input"]))
	}

	e.expect.matcher.input = matcher
	e.expect.matcherLocations["input"] = libCallerLocation(2)
	e.expect.matcherHints["input"] = libMessageMatchArgHint()
	return e
}

// done - MethodExpecterValueArgData.structCode()
type targetFullExpecterWithValueArg struct { // skip:!expect
	expect *targetFullExpect
	target *targetFull
}

// done - MethodExpecterValueArgData.GenerateCode()
func (e *targetFullExpecterWithValueArg) Return(first []Result, second error) { // skip:!expect
	e.expect.returns = targetFullReturn{first: first, second: second}
}

// done - MethodExpecterValueArgData.GenerateCode()
func (e *targetFullExpecterWithValueArg) WithCtx(ctx context.Context) *targetFullExpecterWithValueArg { // skip:!expect
	if e.expect.matcher.ctx != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "WithCtx", e.expect.matcherLocations["ctx"]))
	}

	e.expect.matcher.ctx = libReflectEqualMatcher(ctx)
	e.expect.matcherWants["ctx"] = ctx
	e.expect.matcherMethods["ctx"] = "reflect.DeepEqual"
	e.expect.matcherLocations["ctx"] = libCallerLocation(2)

	return e
}

// done - MethodExpecterValueArgData.GenerateCode()
func (e *targetFullExpecterWithValueArg) WithInput(input string) *targetFullExpecterWithValueArg { // skip:!expect
	if e.expect.matcher.input != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "WithInput", e.expect.matcherLocations["input"]))
	}

	e.expect.matcher.input = libBasicComparisonMatcher(input)
	e.expect.matcherWants["input"] = input
	e.expect.matcherMethods["input"] = "=="
	e.expect.matcherLocations["input"] = libCallerLocation(2)

	return e
}
