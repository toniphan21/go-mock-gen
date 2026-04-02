package mockgen

import "testing"

func Test_MethodExpecterData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodExpecterData
		expected string
	}{
		{
			name: "emit nothing if skip skip expect",
			data: MethodExpecterData{SkipExpect: true},
		},

		{
			name: "emit nothing if no arguments and no returns",
			data: MethodExpecterData{
				TargetMethodStruct:                 "targetMethod",
				TargetMethodReturnStruct:           "targetMethodReturn",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodExpecterStruct:         "targetMethodExpecter",
				TargetMethodExpecterMatchStruct:    "targetMethodExpecterMatch",
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodExpecterValueStruct:    "targetMethodExpecterWithValue",
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				Lib:                                libData(),
			},
			expected: ``,
		},

		{
			name: "emit definition, match and with if there is an argument",
			data: MethodExpecterData{
				TargetMethodStruct:                 "targetMethod",
				TargetMethodReturnStruct:           "targetMethodReturn",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodExpecterStruct:         "targetMethodExpecter",
				TargetMethodExpecterMatchStruct:    "targetMethodExpecterMatch",
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodExpecterValueStruct:    "targetMethodExpecterWithValue",
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				Lib:                                libData(),
				Arguments:                          varInfos("name: name string"),
			},
			expected: `package emitter

type targetMethodExpecter struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecter) Match(matcher func(name string) bool) *targetMethodExpecterMatch {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchByNil(e.target))
	}

	e.expect.match = matcher
	e.expect.matchLocation = libCallerLocation(2)
	return &targetMethodExpecterMatch{expect: e.expect}
}

func (e *targetMethodExpecterWithMatchArg) MatchName(matcher func(name string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchName"))
	}

	e.expect.matcher.name = matcher
	e.expect.matcherLocations["name"] = libCallerLocation(2)
	e.expect.matcherHints["name"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithMatchArg{expect: e.expect, target: e.target}
}

func (e *targetMethodExpecter) With(name string) *targetMethodExpecterWithValue {
	e.WithName(name)
	e.expect.matcherLocations["name"] = libCallerLocation(2)

	return &targetMethodExpecterWithValue{expect: e.expect}
}

func (e *targetMethodExpecterWithValueArg) MatchName(matcher func(name string) bool) *targetMethodExpecterWithValueArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchName"))
	}

	e.expect.matcher.name = matcher
	e.expect.matcherLocations["name"] = libCallerLocation(2)
	e.expect.matcherHints["name"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithValueArg{expect: e.expect, target: e.target}
}
`,
		},

		{
			name: "emit definition, return if there is a return",
			data: MethodExpecterData{
				TargetMethodStruct:                 "targetMethod",
				TargetMethodReturnStruct:           "targetMethodReturn",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodExpecterStruct:         "targetMethodExpecter",
				TargetMethodExpecterMatchStruct:    "targetMethodExpecterMatch",
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodExpecterValueStruct:    "targetMethodExpecterWithValue",
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				Lib:                                libData(),
				Returns:                            varInfos("First: first *time.Time", "Second: second error"),
			},
			expected: `package emitter

import "time"

type targetMethodExpecter struct {
	expect *targetMethodExpect
}

func (e *targetMethodExpecter) Return(first *time.Time, second error) {
	e.expect.returns = targetMethodReturn{First: first, Second: second}
}
`,
		},

		{
			name: "emit definition, return, match and with if there is an argument and return",
			data: MethodExpecterData{
				TargetMethodStruct:                 "targetMethod",
				TargetMethodReturnStruct:           "targetMethodReturn",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodExpecterStruct:         "targetMethodExpecter",
				TargetMethodExpecterMatchStruct:    "targetMethodExpecterMatch",
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodExpecterValueStruct:    "targetMethodExpecterWithValue",
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				Lib:                                libData(),
				Arguments:                          varInfos("val: val map[string]int"),
				Returns:                            varInfos("First: first error"),
			},
			expected: `package emitter

type targetMethodExpecter struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecter) Return(first error) {
	e.expect.returns = targetMethodReturn{First: first}
}

func (e *targetMethodExpecter) Match(matcher func(val map[string]int) bool) *targetMethodExpecterMatch {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchByNil(e.target))
	}

	e.expect.match = matcher
	e.expect.matchLocation = libCallerLocation(2)
	return &targetMethodExpecterMatch{expect: e.expect}
}

func (e *targetMethodExpecterWithMatchArg) MatchVal(matcher func(val map[string]int) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchVal"))
	}

	e.expect.matcher.val = matcher
	e.expect.matcherLocations["val"] = libCallerLocation(2)
	e.expect.matcherHints["val"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithMatchArg{expect: e.expect, target: e.target}
}

func (e *targetMethodExpecter) With(val map[string]int) *targetMethodExpecterWithValue {
	e.WithVal(val)
	e.expect.matcherLocations["val"] = libCallerLocation(2)

	return &targetMethodExpecterWithValue{expect: e.expect}
}

func (e *targetMethodExpecterWithValueArg) MatchVal(matcher func(val map[string]int) bool) *targetMethodExpecterWithValueArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchVal"))
	}

	e.expect.matcher.val = matcher
	e.expect.matcherLocations["val"] = libCallerLocation(2)
	e.expect.matcherHints["val"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithValueArg{expect: e.expect, target: e.target}
}
`,
		},

		{
			name: "emit different receiver name to avoid collision with param name",
			data: MethodExpecterData{
				TargetMethodStruct:                 "targetMethod",
				TargetMethodReturnStruct:           "targetMethodReturn",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodExpecterStruct:         "targetMethodExpecter",
				TargetMethodExpecterMatchStruct:    "targetMethodExpecterMatch",
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodExpecterValueStruct:    "targetMethodExpecterWithValue",
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				Lib:                                libData(),
				Arguments:                          varInfos("e: e string", "e0: e0 string"),
				Returns:                            varInfos("e1: e1 string", "e2: e2 string"),
			},
			expected: `package emitter

type targetMethodExpecter struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e3 *targetMethodExpecter) Return(e1 string, e2 string) {
	e3.expect.returns = targetMethodReturn{e1: e1, e2: e2}
}

func (e3 *targetMethodExpecter) Match(matcher func(e string, e0 string) bool) *targetMethodExpecterMatch {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchByNil(e3.target))
	}

	e3.expect.match = matcher
	e3.expect.matchLocation = libCallerLocation(2)
	return &targetMethodExpecterMatch{expect: e3.expect}
}

func (e3 *targetMethodExpecterWithMatchArg) MatchE(matcher func(e string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE"))
	}

	e3.expect.matcher.e = matcher
	e3.expect.matcherLocations["e"] = libCallerLocation(2)
	e3.expect.matcherHints["e"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithMatchArg{expect: e.expect, target: e.target}
}

func (e3 *targetMethodExpecterWithMatchArg) MatchE0(matcher func(e0 string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE0"))
	}

	e3.expect.matcher.e0 = matcher
	e3.expect.matcherLocations["e0"] = libCallerLocation(2)
	e3.expect.matcherHints["e0"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithMatchArg{expect: e.expect, target: e.target}
}

func (e3 *targetMethodExpecter) With(e string, e0 string) *targetMethodExpecterWithValue {
	e3.WithE(e)
	e3.expect.matcherLocations["e"] = libCallerLocation(2)

	e3.WithE0(e0)
	e3.expect.matcherLocations["e0"] = libCallerLocation(2)

	return &targetMethodExpecterWithValue{expect: e3.expect}
}

func (e3 *targetMethodExpecterWithValueArg) MatchE(matcher func(e string) bool) *targetMethodExpecterWithValueArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE"))
	}

	e3.expect.matcher.e = matcher
	e3.expect.matcherLocations["e"] = libCallerLocation(2)
	e3.expect.matcherHints["e"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithValueArg{expect: e.expect, target: e.target}
}

func (e3 *targetMethodExpecterWithValueArg) MatchE0(matcher func(e0 string) bool) *targetMethodExpecterWithValueArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE0"))
	}

	e3.expect.matcher.e0 = matcher
	e3.expect.matcherLocations["e0"] = libCallerLocation(2)
	e3.expect.matcherHints["e0"] = libMessageMatchArgHint()
	return &targetMethodExpecterWithValueArg{expect: e.expect, target: e.target}
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runEmitterTest(t, &tc.data, tc.expected)
		})
	}
}
