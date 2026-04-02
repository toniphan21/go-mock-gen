package mockgen

import "testing"

func Test_MethodExpecterMatchArgData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodExpecterMatchArgData
		expected string
	}{
		{
			name: "emit nothing if skip skip expect",
			data: MethodExpecterMatchArgData{SkipExpect: true},
		},

		{
			name: "emit nothing if no arguments",
			data: MethodExpecterMatchArgData{},
		},

		{
			name: "emit a definition and a function if there is an argument",
			data: MethodExpecterMatchArgData{
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				Lib:                                libData(),
				Arguments:                          varInfos("name: name string"),
			},
			expected: `package emitter

type targetMethodExpecterWithMatchArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecterWithMatchArg) MatchName(matcher func(name string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchName"))
	}

	if e.expect.matcher.name != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchName", e.expect.matcherLocations["name"]))
	}

	e.expect.matcher.name = matcher
	e.expect.matcherLocations["name"] = libCallerLocation(2)
	e.expect.matcherHints["name"] = libMessageMatchArgHint()
	return e
}
`,
		},

		{
			name: "emit a Return function if Return is not empty",
			data: MethodExpecterMatchArgData{
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodReturnStruct:           "targetMethodReturn",
				Lib:                                libData(),
				Arguments:                          varInfos("name: name string"),
				Returns:                            varInfos("First: first *time.Time", "Second: second error"),
			},
			expected: `package emitter

import "time"

type targetMethodExpecterWithMatchArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecterWithMatchArg) Return(first *time.Time, second error) {
	e.expect.returns = targetMethodReturn{First: first, Second: second}
}

func (e *targetMethodExpecterWithMatchArg) MatchName(matcher func(name string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchName"))
	}

	if e.expect.matcher.name != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchName", e.expect.matcherLocations["name"]))
	}

	e.expect.matcher.name = matcher
	e.expect.matcherLocations["name"] = libCallerLocation(2)
	e.expect.matcherHints["name"] = libMessageMatchArgHint()
	return e
}
`,
		},

		{
			name: "emit different receiver name to avoid collision with param name",
			data: MethodExpecterMatchArgData{
				TargetMethodExpecterMatchArgStruct: "targetMethodExpecterWithMatchArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodReturnStruct:           "targetMethodReturn",
				Lib:                                libData(),
				Arguments:                          varInfos("e: e string", "e0: e0 string"),
				Returns:                            varInfos("e1: e1 string", "e2: e2 string"),
			},
			expected: `package emitter

type targetMethodExpecterWithMatchArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e3 *targetMethodExpecterWithMatchArg) Return(e1 string, e2 string) {
	e3.expect.returns = targetMethodReturn{e1: e1, e2: e2}
}

func (e3 *targetMethodExpecterWithMatchArg) MatchE(matcher func(e string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE"))
	}

	if e3.expect.matcher.e != nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageDuplicateMatchArg(e3.target, "MatchE", e3.expect.matcherLocations["e"]))
	}

	e3.expect.matcher.e = matcher
	e3.expect.matcherLocations["e"] = libCallerLocation(2)
	e3.expect.matcherHints["e"] = libMessageMatchArgHint()
	return e3
}

func (e3 *targetMethodExpecterWithMatchArg) MatchE0(matcher func(e0 string) bool) *targetMethodExpecterWithMatchArg {
	if matcher == nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageMatchArgByNil(e3.target, "MatchE0"))
	}

	if e3.expect.matcher.e0 != nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageDuplicateMatchArg(e3.target, "MatchE0", e3.expect.matcherLocations["e0"]))
	}

	e3.expect.matcher.e0 = matcher
	e3.expect.matcherLocations["e0"] = libCallerLocation(2)
	e3.expect.matcherHints["e0"] = libMessageMatchArgHint()
	return e3
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
