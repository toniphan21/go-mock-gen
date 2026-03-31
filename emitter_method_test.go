package mockgen

import (
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
)

func Test_MethodData_structCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name: "not skip expect",
			expected: `type targetMethod struct {
	Calls        []targetMethodCall
	stub         func()
	stubLocation string
	expects      []*targetMethodExpect
	verified     bool
}
`,
		},

		{
			name:       "skip expect",
			skipExpect: true,
			expected: `type targetMethod struct {
	Calls        []targetMethodCall
	stub         func()
	stubLocation string
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.structCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_callStructCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name: "no arguments, no returns",
			expected: `type targetMethodCall struct {
	Location string
}
`,
		},

		{
			name:      "with arguments, no returns",
			arguments: varInfos("Ctx: ctx context.Context", "Input: input string"),
			expected: `type targetMethodCall struct {
	Location string
	Argument targetMethodArgument
}
`,
		},

		{
			name:      "with arguments, with returns",
			arguments: varInfos("Ctx: ctx context.Context", "Input: input string"),
			returns:   varInfos("First: first string", "Second: second error"),
			expected: `type targetMethodCall struct {
	Location string
	Argument targetMethodArgument
	Return   targetMethodReturn
}
`,
		},

		{
			name:       "with arguments, with returns - skip expect does not affect",
			skipExpect: true,
			arguments:  varInfos("Ctx: ctx context.Context", "Input: input string"),
			returns:    varInfos("First: first string", "Second: second error"),
			expected: `type targetMethodCall struct {
	Location string
	Argument targetMethodArgument
	Return   targetMethodReturn
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.callStructCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_argumentStructCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name:     "no arguments",
			expected: ``,
		},

		{
			name:      "with arguments",
			arguments: varInfos("Ctx: ctx context.Context", "Input: input string"),
			expected: `import "context"

type targetMethodArgument struct {
	Ctx   context.Context
	Input string
}
`,
		},

		{
			name:       "with arguments - skip expect",
			skipExpect: true,
			arguments:  varInfos("Ctx: ctx context.Context", "Input: input string"),
			expected: `import "context"

type targetMethodArgument struct {
	Ctx   context.Context
	Input string
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.argumentStructCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_argumentMatcherStructCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name:     "no arguments",
			expected: ``,
		},

		{
			name:      "with arguments",
			arguments: varInfos("Ctx: ctx context.Context", "Input: input string"),
			expected: `import "context"

type targetMethodArgumentMatcher struct {
	Ctx   func(ctx context.Context) bool
	Input func(input string) bool
}
`,
		},

		{
			name:       "with arguments - skip expect",
			skipExpect: true,
			arguments:  varInfos("Ctx: ctx context.Context", "Input: input string"),
			expected:   ``,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.argumentMatcherStructCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_returnStructCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name:     "no returns",
			expected: ``,
		},

		{
			name:    "has returns",
			returns: varInfos("First: first string", "Second: second error"),
			expected: `type targetMethodReturn struct {
	First  string
	Second error
}
`,
		},

		{
			name:       "has returns, skip expect does not affect",
			skipExpect: true,
			returns:    varInfos("First: first string", "Second: second error"),
			expected: `type targetMethodReturn struct {
	First  string
	Second error
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.returnStructCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_expectStructCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		returns    []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name:       "skip expect",
			skipExpect: true,
			expected:   ``,
		},

		{
			name:    "no arguments",
			returns: varInfos("First: first string", "Second: second string"),
			expected: `import "testing"

type targetMethodExpect struct {
	returns  targetMethodReturn
	location string
	index    int
	tb       testing.TB
}
`,
		},

		{
			name:      "no returns",
			arguments: varInfos("Input: input string"),
			expected: `import "testing"

type targetMethodExpect struct {
	match            func(input string) bool
	matchLocation    string
	matcher          *targetMethodArgumentMatcher
	matcherWants     map[string]any
	matcherMethods   map[string]string
	matcherHints     map[string]string
	matcherLocations map[string]string
	location         string
	index            int
	tb               testing.TB
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.expectStructCode()
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func runMethodDataFunc(arguments []VarInfo, returns []VarInfo, skipExpect bool, fn func(data MethodData) jen.Code) string {
	data := MethodData{
		TargetMethodStruct:                "targetMethod",
		TargetMethodCallStruct:            "targetMethodCall",
		TargetMethodArgumentStruct:        "targetMethodArgument",
		TargetMethodArgumentMatcherStruct: "targetMethodArgumentMatcher",
		TargetMethodReturnStruct:          "targetMethodReturn",
		TargetMethodExpectStruct:          "targetMethodExpect",
		Lib:                               libData(),
		Arguments:                         arguments,
		Returns:                           returns,
		SkipExpect:                        skipExpect,
	}

	code := fn(data)
	if code == nil {
		return ""
	}

	jf := jen.NewFile("emitter")
	jf.Add(code)

	return strings.ReplaceAll(jf.GoString(), "package emitter\n\n", "")
}
