package mockgen

import (
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"nhatp.com/go/gen-lib/gentest"
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

func Test_MethodData_methodNameFuncCode(t *testing.T) {
	out := runMethodDataFunc(nil, nil, false, func(data MethodData) jen.Code {
		return data.methodNameFuncCode("x")
	})

	expected := "func (x *targetMethod) methodName() string {\n\treturn \"Method\"\n}\n"

	assert.Equal(t, expected, out)
}

func Test_MethodData_interfaceNameFuncCode(t *testing.T) {
	out := runMethodDataFunc(nil, nil, false, func(data MethodData) jen.Code {
		return data.interfaceNameFuncCode("x")
	})

	expected := "func (x *targetMethod) interfaceName() string {\n\treturn \"Target\"\n}\n"

	assert.Equal(t, expected, out)
}

func Test_MethodData_fatalFuncCode(t *testing.T) {
	out := runMethodDataFunc(nil, nil, false, func(data MethodData) jen.Code {
		return data.fatalFuncCode("x")
	})

	expected := `func (x *targetMethod) fatal(index int, msg string) {
	x.verified = true
	x.expects[index].tb.Helper()
	x.expects[index].tb.Fatal(msg)
}
`

	assert.Equal(t, expected, out)

	out = runMethodDataFunc(nil, nil, true, func(data MethodData) jen.Code {
		return data.fatalFuncCode("x")
	})

	expected = "func (x *targetMethod) fatal(index int, msg string) {}\n"

	assert.Equal(t, expected, out)
}

func Test_MethodData_panicFuncCode(t *testing.T) {
	out := runMethodDataFunc(nil, nil, false, func(data MethodData) jen.Code {
		return data.panicFuncCode("x")
	})

	expected := `func (x *targetMethod) panic(msg string) {
	x.verified = true
	panic(msg)
}
`

	assert.Equal(t, expected, out)

	out = runMethodDataFunc(nil, nil, true, func(data MethodData) jen.Code {
		return data.panicFuncCode("x")
	})

	expected = `func (x *targetMethod) panic(msg string) {
	panic(msg)
}
`

	assert.Equal(t, expected, out)
}

func Test_MethodData_buildCallHistoryFuncCode(t *testing.T) {
	cases := []struct {
		name       string
		arguments  []VarInfo
		skipExpect bool
		expected   string
	}{
		{
			name:       "no arguments, skip expect",
			skipExpect: true,
			expected: `import "strings"

func (x *targetMethod) buildCallHistory(sb *strings.Builder, header string) {}
`,
		},

		{
			name: "no arguments",
			expected: `import (
	"fmt"
	"strings"
)

func (x *targetMethod) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(x.Calls) != 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, v := range x.Calls {
		a := []any{}
		libMessageCallHistory(sb, i, x.expects[i].location, v.Location, a)
	}
}
`,
		},

		{
			name:      "with arguments",
			arguments: varInfos("Ctx: ctx context.Context", "ID: id int"),
			expected: `import (
	"fmt"
	"strings"
)

func (x *targetMethod) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(x.Calls) != 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, v := range x.Calls {
		a := []any{"ctx", v.Argument.Ctx, "id", v.Argument.ID}
		libMessageCallHistory(sb, i, x.expects[i].location, v.Location, a)
	}
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, nil, tc.skipExpect, func(data MethodData) jen.Code {
				return data.buildCallHistoryFuncCode("x")
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_invokeStubFuncCode(t *testing.T) {
	cases := []struct {
		name      string
		arguments []VarInfo
		returns   []VarInfo
		expected  string
	}{
		{
			name: "no arguments, no returns",
			expected: `func (x *targetMethod) invokeStub() {
	x.stub()
	x.capture()
}
`,
		},

		{
			name:    "no arguments, with 1 return",
			returns: varInfos("First: first error"),
			expected: `func (x *targetMethod) invokeStub() error {
	v0 := x.stub()
	return x.capture(targetMethodReturn{First: v0})
}
`,
		},

		{
			name: "no arguments, with 2 returns",
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeStub() (first string, err error) {
	v0, v1 := x.stub()
	return x.capture(targetMethodReturn{First: v0, Error: v1})
}
`,
		},

		{
			name:      "with arguments, no return",
			arguments: varInfos("Input: input string"),
			expected: `func (x *targetMethod) invokeStub(input string) {
	x.stub(input)
	x.capture(targetMethodArgument{Input: input})
}
`,
		},

		{
			name:      "with arguments, with 1 return",
			arguments: varInfos("Input: input string"),
			returns:   varInfos("First: first error"),
			expected: `func (x *targetMethod) invokeStub(input string) error {
	v0 := x.stub(input)
	return x.capture(targetMethodArgument{Input: input}, targetMethodReturn{First: v0})
}
`,
		},

		{
			name:      "with arguments, with 2 returns",
			arguments: varInfos("Input: input string"),
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeStub(input string) (first string, err error) {
	v0, v1 := x.stub(input)
	return x.capture(targetMethodArgument{Input: input}, targetMethodReturn{First: v0, Error: v1})
}
`,
		},

		{
			name:      "with name collision",
			arguments: varInfos("Input: v0 string"),
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "v1", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "v2", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeStub(v0 string) (v1 string, v2 error) {
	v3, v4 := x.stub(v0)
	return x.capture(targetMethodArgument{Input: v0}, targetMethodReturn{First: v3, Error: v4})
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, true, func(data MethodData) jen.Code {
				return data.invokeStubFuncCode("x")
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_invokeExpectFuncCode(t *testing.T) {
	cases := []struct {
		name       string
		skipExpect bool
		arguments  []VarInfo
		returns    []VarInfo
		expected   string
	}{
		{
			name:       "skip expect",
			skipExpect: true,
			expected:   ``,
		},

		{
			name: "no arguments, no returns",
			expected: `func (x *targetMethod) invokeExpect() {
	v0 := []any{}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	x.capture()
}
`,
		},

		{
			name:    "no arguments, with 1 return",
			returns: varInfos("First: first error"),
			expected: `func (x *targetMethod) invokeExpect() error {
	v0 := []any{}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	v2 := x.expects[v1]
	return x.capture(v2.returns)
}
`,
		},

		{
			name: "no arguments, with 2 returns",
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeExpect() (first string, err error) {
	v0 := []any{}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	v2 := x.expects[v1]
	return x.capture(v2.returns)
}
`,
		},

		{
			name:      "with arguments, no return",
			arguments: varInfos("Input: input string"),
			expected: `func (x *targetMethod) invokeExpect(input string) {
	v0 := []any{"input", input}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	v2 := x.expects[v1]
	if v2.match != nil && !v2.match(input) {
		v2.tb.Helper()
		x.fatal(v1, libMessageMatchFail(x, v2.matchLocation, v1, v0))
	}

	v2.tb.Helper()
	libMatchArgument(x, v1, "input", input, v2.matcher.input, v2.matcherWants, v2.matcherMethods, v2.matcherHints, v2.tb, v2.matcherLocations["input"])

	x.capture(targetMethodArgument{Input: input})
}
`,
		},

		{
			name:      "with arguments, with 1 return",
			arguments: varInfos("Ctx: ctx context.Context", "ID: id int"),
			returns:   varInfos("First: first error"),
			expected: `import "context"

func (x *targetMethod) invokeExpect(ctx context.Context, id int) error {
	v0 := []any{"ctx", ctx, "id", id}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	v2 := x.expects[v1]
	if v2.match != nil && !v2.match(ctx, id) {
		v2.tb.Helper()
		x.fatal(v1, libMessageMatchFail(x, v2.matchLocation, v1, v0))
	}

	v2.tb.Helper()
	libMatchArgument(x, v1, "ctx", ctx, v2.matcher.ctx, v2.matcherWants, v2.matcherMethods, v2.matcherHints, v2.tb, v2.matcherLocations["ctx"])
	libMatchArgument(x, v1, "id", id, v2.matcher.id, v2.matcherWants, v2.matcherMethods, v2.matcherHints, v2.tb, v2.matcherLocations["id"])

	return x.capture(targetMethodArgument{Ctx: ctx, ID: id}, v2.returns)
}
`,
		},

		{
			name:      "with arguments, with 2 returns",
			arguments: varInfos("Input: input string"),
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeExpect(input string) (first string, err error) {
	v0 := []any{"input", input}

	v1 := len(x.Calls)
	if v1 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v1+1, v0))
	}

	v2 := x.expects[v1]
	if v2.match != nil && !v2.match(input) {
		v2.tb.Helper()
		x.fatal(v1, libMessageMatchFail(x, v2.matchLocation, v1, v0))
	}

	v2.tb.Helper()
	libMatchArgument(x, v1, "input", input, v2.matcher.input, v2.matcherWants, v2.matcherMethods, v2.matcherHints, v2.tb, v2.matcherLocations["input"])

	return x.capture(targetMethodArgument{Input: input}, v2.returns)
}
`,
		},

		{
			name:      "with name collision",
			arguments: varInfos("Input: v0 string"),
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "v1", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "v2", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) invokeExpect(v0 string) (v1 string, v2 error) {
	v3 := []any{"v0", v0}

	v4 := len(x.Calls)
	if v4 >= len(x.expects) {
		x.panic(libMessageTooManyCalls(x, len(x.expects), v4+1, v3))
	}

	v5 := x.expects[v4]
	if v5.match != nil && !v5.match(v0) {
		v5.tb.Helper()
		x.fatal(v4, libMessageMatchFail(x, v5.matchLocation, v4, v3))
	}

	v5.tb.Helper()
	libMatchArgument(x, v4, "v0", v0, v5.matcher.v0, v5.matcherWants, v5.matcherMethods, v5.matcherHints, v5.tb, v5.matcherLocations["v0"])

	return x.capture(targetMethodArgument{Input: v0}, v5.returns)
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, tc.skipExpect, func(data MethodData) jen.Code {
				return data.invokeExpectFuncCode("x")
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_captureFuncCode(t *testing.T) {
	cases := []struct {
		name      string
		arguments []VarInfo
		returns   []VarInfo
		expected  string
	}{
		{
			name: "no arguments, no returns",
			expected: `func (x *targetMethod) capture() {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4)})
}
`,
		},

		{
			name:    "no arguments, with 1 return",
			returns: varInfos("First: first error"),
			expected: `func (x *targetMethod) capture(returns targetMethodReturn) error {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4), Return: returns})
	return returns.First
}
`,
		},

		{
			name: "no arguments, with 2 returns",
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) capture(returns targetMethodReturn) (first string, err error) {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4), Return: returns})
	return returns.First, returns.Error
}
`,
		},

		{
			name:      "with arguments, no return",
			arguments: varInfos("Input: input string"),
			expected: `func (x *targetMethod) capture(args targetMethodArgument) {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4), Argument: args})
}
`,
		},

		{
			name:      "with arguments, with 1 return",
			arguments: varInfos("Input: input string"),
			returns:   varInfos("First: first error"),
			expected: `func (x *targetMethod) capture(args targetMethodArgument, returns targetMethodReturn) error {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4), Argument: args, Return: returns})
	return returns.First
}
`,
		},

		{
			name:      "with arguments, with 2 returns",
			arguments: varInfos("Input: input string"),
			returns: []VarInfo{
				{Name: "first", Field: "First", OriginalName: "first", Type: gentest.Type("string")},
				{Name: "second", Field: "Error", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `func (x *targetMethod) capture(args targetMethodArgument, returns targetMethodReturn) (first string, err error) {
	x.Calls = append(x.Calls, targetMethodCall{Location: libCallerLocation(4), Argument: args, Return: returns})
	return returns.First, returns.Error
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := runMethodDataFunc(tc.arguments, tc.returns, true, func(data MethodData) jen.Code {
				return data.captureFuncCode("x")
			})

			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_MethodData_verifyFuncCode(t *testing.T) {
	out := runMethodDataFunc(nil, nil, false, func(data MethodData) jen.Code {
		return data.verifyFuncCode("x")
	})

	expected := `func (x *targetMethod) verify(index int) {
	if !x.verified && index >= len(x.Calls) {
		x.expects[index].tb.Helper()
		x.expects[index].tb.Fatal(libMessageExpectButNotCalled(x, len(x.expects), len(x.Calls), index))
	}
}
`

	assert.Equal(t, expected, out)

	out = runMethodDataFunc(nil, nil, true, func(data MethodData) jen.Code {
		return data.verifyFuncCode("x")
	})

	expected = ""

	assert.Equal(t, expected, out)
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
		Struct:                "targetMethod",
		CallStruct:            "targetMethodCall",
		ArgumentStruct:        "targetMethodArgument",
		ArgumentMatcherStruct: "targetMethodArgumentMatcher",
		ReturnStruct:          "targetMethodReturn",
		ExpectStruct:          "targetMethodExpect",
		Interface:             "Target",
		Name:                  "Method",
		Lib:                   libData(),
		Arguments:             arguments,
		Returns:               returns,
		SkipExpect:            skipExpect,
	}

	code := fn(data)
	if code == nil {
		return ""
	}

	jf := jen.NewFile("emitter")
	jf.Add(code)

	return strings.ReplaceAll(jf.GoString(), "package emitter\n\n", "")
}

func Test_MethodData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodData
		expected string
	}{
		{
			name: "no arguments, no returns",
			data: MethodData{
				Struct:                "targetMethod",
				CallStruct:            "targetMethodCall",
				ArgumentStruct:        "targetMethodArgument",
				ArgumentMatcherStruct: "targetMethodArgumentMatcher",
				ReturnStruct:          "targetMethodReturn",
				ExpectStruct:          "targetMethodExpect",
				Interface:             "Target",
				Name:                  "Method",
				Lib:                   libData(),
			},
			expected: `package emitter

import (
	"fmt"
	"strings"
	"testing"
)

type targetMethod struct {
	Calls        []targetMethodCall
	stub         func()
	stubLocation string
	expects      []*targetMethodExpect
	verified     bool
}

func (m *targetMethod) methodName() string {
	return "Method"
}

func (m *targetMethod) interfaceName() string {
	return "Target"
}

func (m *targetMethod) fatal(index int, msg string) {
	m.verified = true
	m.expects[index].tb.Helper()
	m.expects[index].tb.Fatal(msg)
}

func (m *targetMethod) panic(msg string) {
	m.verified = true
	panic(msg)
}

func (m *targetMethod) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(m.Calls) != 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, v := range m.Calls {
		a := []any{}
		libMessageCallHistory(sb, i, m.expects[i].location, v.Location, a)
	}
}

func (m *targetMethod) invokeStub() {
	m.stub()
	m.capture()
}

func (m *targetMethod) invokeExpect() {
	v0 := []any{}

	v1 := len(m.Calls)
	if v1 >= len(m.expects) {
		m.panic(libMessageTooManyCalls(m, len(m.expects), v1+1, v0))
	}

	m.capture()
}

func (m *targetMethod) capture() {
	m.Calls = append(m.Calls, targetMethodCall{Location: libCallerLocation(4)})
}

func (m *targetMethod) verify(index int) {
	if !m.verified && index >= len(m.Calls) {
		m.expects[index].tb.Helper()
		m.expects[index].tb.Fatal(libMessageExpectButNotCalled(m, len(m.expects), len(m.Calls), index))
	}
}

type targetMethodCall struct {
	Location string
}

type targetMethodExpect struct {
	location string
	index    int
	tb       testing.TB
}
`,
		},

		{
			name: "no arguments, no returns, skip expect",
			data: MethodData{
				Struct:                "targetMethod",
				CallStruct:            "targetMethodCall",
				ArgumentStruct:        "targetMethodArgument",
				ArgumentMatcherStruct: "targetMethodArgumentMatcher",
				ReturnStruct:          "targetMethodReturn",
				ExpectStruct:          "targetMethodExpect",
				Interface:             "Target",
				Name:                  "Method",
				Lib:                   libData(),
				SkipExpect:            true,
			},
			expected: `package emitter

import "strings"

type targetMethod struct {
	Calls        []targetMethodCall
	stub         func()
	stubLocation string
}

func (m *targetMethod) methodName() string {
	return "Method"
}

func (m *targetMethod) interfaceName() string {
	return "Target"
}

func (m *targetMethod) fatal(index int, msg string) {}

func (m *targetMethod) panic(msg string) {
	panic(msg)
}

func (m *targetMethod) buildCallHistory(sb *strings.Builder, header string) {}

func (m *targetMethod) invokeStub() {
	m.stub()
	m.capture()
}

func (m *targetMethod) capture() {
	m.Calls = append(m.Calls, targetMethodCall{Location: libCallerLocation(4)})
}

type targetMethodCall struct {
	Location string
}
`,
		},

		{
			name: "handle name collision",
			data: MethodData{
				Struct:                "targetMethod",
				CallStruct:            "targetMethodCall",
				ArgumentStruct:        "targetMethodArgument",
				ArgumentMatcherStruct: "targetMethodArgumentMatcher",
				ReturnStruct:          "targetMethodReturn",
				ExpectStruct:          "targetMethodExpect",
				Interface:             "Target",
				Name:                  "Method",
				Lib:                   libData(),
				Arguments:             varInfos("Input: m string"),
				Returns: []VarInfo{
					{Name: "first", Field: "First", OriginalName: "m0", Type: gentest.Type("string")},
					{Name: "second", Field: "Error", OriginalName: "m1", Type: gentest.Type("error")},
				},
			},
			expected: `package emitter

import (
	"fmt"
	"strings"
	"testing"
)

type targetMethod struct {
	Calls        []targetMethodCall
	stub         func(m string) (m0 string, m1 error)
	stubLocation string
	expects      []*targetMethodExpect
	verified     bool
}

func (m2 *targetMethod) methodName() string {
	return "Method"
}

func (m2 *targetMethod) interfaceName() string {
	return "Target"
}

func (m2 *targetMethod) fatal(index int, msg string) {
	m2.verified = true
	m2.expects[index].tb.Helper()
	m2.expects[index].tb.Fatal(msg)
}

func (m2 *targetMethod) panic(msg string) {
	m2.verified = true
	panic(msg)
}

func (m2 *targetMethod) buildCallHistory(sb *strings.Builder, header string) {
	if header != "" && len(m2.Calls) != 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", header))
	}

	for i, v := range m2.Calls {
		a := []any{"m", v.Argument.Input}
		libMessageCallHistory(sb, i, m2.expects[i].location, v.Location, a)
	}
}

func (m2 *targetMethod) invokeStub(m string) (m0 string, m1 error) {
	v0, v1 := m2.stub(m)
	return m2.capture(targetMethodArgument{Input: m}, targetMethodReturn{First: v0, Error: v1})
}

func (m2 *targetMethod) invokeExpect(m string) (m0 string, m1 error) {
	v0 := []any{"m", m}

	v1 := len(m2.Calls)
	if v1 >= len(m2.expects) {
		m2.panic(libMessageTooManyCalls(m2, len(m2.expects), v1+1, v0))
	}

	v2 := m2.expects[v1]
	if v2.match != nil && !v2.match(m) {
		v2.tb.Helper()
		m2.fatal(v1, libMessageMatchFail(m2, v2.matchLocation, v1, v0))
	}

	v2.tb.Helper()
	libMatchArgument(m2, v1, "m", m, v2.matcher.m, v2.matcherWants, v2.matcherMethods, v2.matcherHints, v2.tb, v2.matcherLocations["m"])

	return m2.capture(targetMethodArgument{Input: m}, v2.returns)
}

func (m2 *targetMethod) capture(args targetMethodArgument, returns targetMethodReturn) (m0 string, m1 error) {
	m2.Calls = append(m2.Calls, targetMethodCall{Location: libCallerLocation(4), Argument: args, Return: returns})
	return returns.First, returns.Error
}

func (m2 *targetMethod) verify(index int) {
	if !m2.verified && index >= len(m2.Calls) {
		m2.expects[index].tb.Helper()
		m2.expects[index].tb.Fatal(libMessageExpectButNotCalled(m2, len(m2.expects), len(m2.Calls), index))
	}
}

type targetMethodCall struct {
	Location string
	Argument targetMethodArgument
	Return   targetMethodReturn
}

type targetMethodArgument struct {
	Input string
}

type targetMethodArgumentMatcher struct {
	Input func(m string) bool
}

type targetMethodReturn struct {
	First string
	Error error
}

type targetMethodExpect struct {
	match            func(m string) bool
	matchLocation    string
	matcher          *targetMethodArgumentMatcher
	matcherWants     map[string]any
	matcherMethods   map[string]string
	matcherHints     map[string]string
	matcherLocations map[string]string
	returns          targetMethodReturn
	location         string
	index            int
	tb               testing.TB
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
