package mockgen

import (
	"testing"
)

func Test_MethodExpecterValueArgData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodExpecterValueArgData
		expected string
	}{
		{
			name: "emit nothing if skip skip expect",
			data: MethodExpecterValueArgData{SkipExpect: true},
		},

		{
			name: "emit nothing if no arguments",
			data: MethodExpecterValueArgData{},
		},

		{
			name: "emit a definition and a function if there is an argument",
			data: MethodExpecterValueArgData{
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				Lib:                                libData(),
				Arguments:                          varInfos("name: name string"),
			},
			expected: `package emitter

type targetMethodExpecterWithValueArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecterWithValueArg) WithName(name string) *targetMethodExpecterWithValueArg {
	if e.expect.matcher.name != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "WithName", e.expect.matcherLocations["name"]))
	}

	e.expect.matcher.name = libBasicComparisonMatcher(name)
	e.expect.matcherWants["name"] = name
	e.expect.matcherMethods["name"] = "=="
	e.expect.matcherLocations["name"] = libCallerLocation(2)

	return e
}
`,
		},

		{
			name: "emit code use reflect.DeepEqual if the type is not comparable",
			data: MethodExpecterValueArgData{
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				Lib:                                libData(),
				Arguments:                          varInfos("val: val map[string]int"),
			},
			expected: `package emitter

type targetMethodExpecterWithValueArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecterWithValueArg) WithVal(val map[string]int) *targetMethodExpecterWithValueArg {
	if e.expect.matcher.val != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "WithVal", e.expect.matcherLocations["val"]))
	}

	e.expect.matcher.val = libReflectEqualMatcher(val)
	e.expect.matcherWants["val"] = val
	e.expect.matcherMethods["val"] = "reflect.DeepEqual"
	e.expect.matcherLocations["val"] = libCallerLocation(2)

	return e
}
`,
		},

		{
			name: "emit a Return function if Return is not empty",
			data: MethodExpecterValueArgData{
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodReturnStruct:           "targetMethodReturn",
				Lib:                                libData(),
				Arguments:                          varInfos("name: name string"),
				Returns:                            varInfos("First: first *time.Time", "Second: second error"),
			},
			expected: `package emitter

import "time"

type targetMethodExpecterWithValueArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e *targetMethodExpecterWithValueArg) Return(first *time.Time, second error) {
	e.expect.returns = targetMethodReturn{First: first, Second: second}
}

func (e *targetMethodExpecterWithValueArg) WithName(name string) *targetMethodExpecterWithValueArg {
	if e.expect.matcher.name != nil {
		e.expect.tb.Helper()
		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "WithName", e.expect.matcherLocations["name"]))
	}

	e.expect.matcher.name = libBasicComparisonMatcher(name)
	e.expect.matcherWants["name"] = name
	e.expect.matcherMethods["name"] = "=="
	e.expect.matcherLocations["name"] = libCallerLocation(2)

	return e
}
`,
		},

		{
			name: "emit different receiver name to avoid collision with param name",
			data: MethodExpecterValueArgData{
				TargetMethodExpecterValueArgStruct: "targetMethodExpecterWithValueArg",
				TargetMethodStruct:                 "targetMethod",
				TargetMethodExpectStruct:           "targetMethodExpect",
				TargetMethodReturnStruct:           "targetMethodReturn",
				Lib:                                libData(),
				Arguments:                          varInfos("e: e string", "e0: e0 string"),
				Returns:                            varInfos("e1: e1 string", "e2: e2 string"),
			},
			expected: `package emitter

type targetMethodExpecterWithValueArg struct {
	expect *targetMethodExpect
	target *targetMethod
}

func (e3 *targetMethodExpecterWithValueArg) Return(e1 string, e2 string) {
	e3.expect.returns = targetMethodReturn{e1: e1, e2: e2}
}

func (e3 *targetMethodExpecterWithValueArg) WithE(e string) *targetMethodExpecterWithValueArg {
	if e3.expect.matcher.e != nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageDuplicateMatchArg(e3.target, "WithE", e3.expect.matcherLocations["e"]))
	}

	e3.expect.matcher.e = libBasicComparisonMatcher(e)
	e3.expect.matcherWants["e"] = e
	e3.expect.matcherMethods["e"] = "=="
	e3.expect.matcherLocations["e"] = libCallerLocation(2)

	return e3
}

func (e3 *targetMethodExpecterWithValueArg) WithE0(e0 string) *targetMethodExpecterWithValueArg {
	if e3.expect.matcher.e0 != nil {
		e3.expect.tb.Helper()
		e3.target.fatal(e3.expect.index, libMessageDuplicateMatchArg(e3.target, "WithE0", e3.expect.matcherLocations["e0"]))
	}

	e3.expect.matcher.e0 = libBasicComparisonMatcher(e0)
	e3.expect.matcherWants["e0"] = e0
	e3.expect.matcherMethods["e0"] = "=="
	e3.expect.matcherLocations["e0"] = libCallerLocation(2)

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
