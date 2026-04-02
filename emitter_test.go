package mockgen

import (
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhatp.com/go/gen-lib/gentest"
)

type emitter interface {
	GenerateCode() []jen.Code
}

func runEmitterTest[T emitter](t *testing.T, em T, expected string) {
	code := em.GenerateCode()
	if expected == "" {
		assert.Nil(t, code)
		return
	}

	assert.NotNil(t, code)

	jf := jen.NewFile("emitter")
	for _, v := range code {
		jf.Add(v)
	}
	out := jf.GoString()
	assert.Equal(t, expected, out)
}

func varInfos(args ...string) []VarInfo {
	var result []VarInfo
	for _, arg := range args {
		v := strings.Split(arg, ":")
		if len(v) != 2 {
			panic("invalid VarInfo: " + arg)
		}

		p := strings.Split(strings.TrimSpace(v[1]), " ")
		if len(p) != 2 {
			panic("invalid VarInfo: " + arg)
		}
		t := gentest.Param(p[0], p[1])
		result = append(result, VarInfo{Field: v[0], Name: p[0], Type: t.Type()})
	}
	return result
}

func Test_targetMethodSignatureString(t *testing.T) {
	cases := []struct {
		name     string
		method   MethodInfo
		expected string
	}{
		{
			name:     "empty",
			method:   MethodInfo{},
			expected: "()",
		},

		{
			name: "has 1 argument",
			method: MethodInfo{
				Arguments: []VarInfo{
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
			},
			expected: "(input string)",
		},

		{
			name: "has 2 arguments",
			method: MethodInfo{
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
			},
			expected: "(ctx Context, input string)",
		},

		{
			name: "has 1 return",
			method: MethodInfo{
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("error")},
				},
			},
			expected: "() error",
		},

		{
			name: "has 1 named return",
			method: MethodInfo{
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
				},
			},
			expected: "() (err error)",
		},

		{
			name: "has 2 returns",
			method: MethodInfo{
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("string")},
					{Name: "second", Field: "second", OriginalName: "", Type: gentest.Type("error")},
				},
			},
			expected: "() (string, error)",
		},

		{
			name: "has 1 argument, 1 return",
			method: MethodInfo{
				Arguments: []VarInfo{
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("error")},
				},
			},
			expected: "(input string) error",
		},

		{
			name: "has 2 arguments, 1 named return",
			method: MethodInfo{
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
				},
			},
			expected: "(ctx Context, input string) (err error)",
		},

		{
			name: "has 2 arguments, 2 returns",
			method: MethodInfo{
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("string")},
					{Name: "second", Field: "second", OriginalName: "", Type: gentest.Type("error")},
				},
			},
			expected: "(ctx Context, input string) (string, error)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := targetMethodSignatureString(tc.method)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func Test_targetMethodSignature(t *testing.T) {
	cases := []struct {
		name      string
		arguments []VarInfo
		returns   []VarInfo
		expected  string
	}{
		{
			name:     "empty",
			expected: "var out func()",
		},

		{
			name: "has 1 argument",
			arguments: []VarInfo{
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: "var out func(input string)",
		},

		{
			name: "has a unnamed argument",
			arguments: []VarInfo{
				{Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: "var out func(string)",
		},

		{
			name: "has 2 arguments",
			arguments: []VarInfo{
				{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: `import "context"` + "\n\n" + `var out func(ctx context.Context, input string)`,
		},

		{
			name: "has 1 return",
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("error")},
			},
			expected: "var out func() error",
		},

		{
			name: "has 1 named return",
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: "var out func() (err error)",
		},

		{
			name: "has 2 returns",
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("string")},
				{Name: "second", Field: "second", OriginalName: "", Type: gentest.Type("error")},
			},
			expected: "var out func() (string, error)",
		},

		{
			name: "has 1 argument, 1 return",
			arguments: []VarInfo{
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("error")},
			},
			expected: "var out func(input string) error",
		},

		{
			name: "has 2 arguments, 1 named return",
			arguments: []VarInfo{
				{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
			},
			expected: `import "context"` + "\n\n" + "var out func(ctx context.Context, input string) (err error)",
		},

		{
			name: "has 2 arguments, 2 returns",
			arguments: []VarInfo{
				{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			returns: []VarInfo{
				{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("string")},
				{Name: "second", Field: "second", OriginalName: "", Type: gentest.Type("error")},
			},
			expected: `import "context"` + "\n\n" + "var out func(ctx context.Context, input string) (string, error)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code := targetMethodSignature(tc.arguments, tc.returns)
			require.NotNil(t, code)

			jf := jen.NewFile("emitter")
			jf.Var().Id("out").Add(code)
			out := strings.ReplaceAll(jf.GoString(), "package emitter\n\n", "")
			assert.Equal(t, tc.expected+"\n", out)
		})
	}
}

func Test_targetMethodMatcherSignature(t *testing.T) {
	cases := []struct {
		name      string
		arguments []VarInfo
		expected  string
	}{
		{
			name:      "empty",
			arguments: nil,
			expected:  "",
		},

		{
			name: "has 1 argument",
			arguments: []VarInfo{
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: "var out func(input string) bool",
		},

		{
			name: "has unnamed argument",
			arguments: []VarInfo{
				{Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: "var out func(string) bool",
		},

		{
			name: "has 2 arguments",
			arguments: []VarInfo{
				{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
				{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
			},
			expected: `import "context"` + "\n\n" + `var out func(ctx context.Context, input string) bool`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code := targetMethodMatcherSignature(tc.arguments...)
			if tc.expected == "" {
				assert.Nil(t, code)
				return
			}

			require.NotNil(t, code)

			jf := jen.NewFile("emitter")
			jf.Var().Id("out").Add(code)
			out := strings.ReplaceAll(jf.GoString(), "package emitter\n\n", "")
			assert.Equal(t, tc.expected+"\n", out)
		})
	}
}
