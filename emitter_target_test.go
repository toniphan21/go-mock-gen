package mockgen

import (
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"nhatp.com/go/gen-lib/gentest"
)

func Test_TargetData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     TargetData
		expected string
	}{
		{
			name:     "return nil if there is no method",
			data:     TargetData{},
			expected: ``,
		},

		{
			name: "skip constructor if it is empty",
			data: TargetData{
				Interface:        "Target",
				Struct:           "target",
				Constructor:      "",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				ExpecterStruct:   "targetExpecter",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:   "Method",
						Struct: "targetMethod",
					},
				},
			},
			expected: `package emitter

type targetTestDouble struct {
	location string
	Method   *targetMethod
}

type target struct {
	td *targetTestDouble
}

func (m *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

func (m *target) EXPECT() *targetExpecter {
	return &targetExpecter{target: m}
}

func (m *target) Method() {
	v0, v1, v2 := "Target", "Method", "()"
	v3 := []any{}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub()
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect()
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name: "strip the expected if it is skipped",
			data: TargetData{
				Interface:        "Target",
				Struct:           "target",
				Constructor:      "testTarget",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				ExpecterStruct:   "targetExpecter",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:   "Method",
						Struct: "targetMethod",
					},
				},
			},
			expected: `package emitter

func testTarget() *target {
	return &target{td: &targetTestDouble{location: libCallerLocation(2)}}
}

type targetTestDouble struct {
	location string
	Method   *targetMethod
}

type target struct {
	td *targetTestDouble
}

func (m *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

func (m *target) EXPECT() *targetExpecter {
	return &targetExpecter{target: m}
}

func (m *target) Method() {
	v0, v1, v2 := "Target", "Method", "()"
	v3 := []any{}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub()
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect()
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name: "test double will choose other name if there is a method named 'location'",
			data: TargetData{
				Interface:        "Target",
				Struct:           "target",
				Constructor:      "testTarget",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				ExpecterStruct:   "targetExpecter",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:   "location",
						Struct: "targetlocation",
					},
				},
			},
			expected: `package emitter

func testTarget() *target {
	return &target{td: &targetTestDouble{location0: libCallerLocation(2)}}
}

type targetTestDouble struct {
	location0 string
	location  *targetlocation
}

type target struct {
	td *targetTestDouble
}

func (m *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

func (m *target) EXPECT() *targetExpecter {
	return &targetExpecter{target: m}
}

func (m *target) location() {
	v0, v1, v2 := "Target", "location", "()"
	v3 := []any{}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.location; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub()
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect()
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location0, v3))
}
`,
		},

		{
			name: "target can avoid receiver name collision",
			data: TargetData{
				Interface:        "Target",
				Struct:           "target",
				Constructor:      "testTarget",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				ExpecterStruct:   "targetExpecter",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:      "Method",
						Struct:    "targetMocation",
						Arguments: varInfos("m: m string"),
					},
					{
						Name:      "Awesome",
						Struct:    "targetAwesome",
						Arguments: varInfos("m0: m0 int"),
						Returns: []VarInfo{
							{Name: "m1", OriginalName: "m1", Field: "m1", Type: gentest.Type("error")},
						},
					},
				},
			},
			expected: `package emitter

func testTarget() *target {
	return &target{td: &targetTestDouble{location: libCallerLocation(2)}}
}

type targetTestDouble struct {
	location string
	Method   *targetMocation
	Awesome  *targetAwesome
}

type target struct {
	td *targetTestDouble
}

func (m2 *target) STUB() *targetStubber {
	return &targetStubber{target: m}
}

func (m2 *target) EXPECT() *targetExpecter {
	return &targetExpecter{target: m}
}

func (m2 *target) Method(m string) {
	v0, v1, v2 := "Target", "Method", "(m string)"
	v3 := []any{"m", m}

	if m2.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m2.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub(m)
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect(m)
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m2.td.location, v3))
}

func (m2 *target) Awesome(m0 int) (m1 error) {
	v0, v1, v2 := "Target", "Awesome", "(m0 int) (m1 error)"
	v3 := []any{"m0", m0}

	if m2.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m2.td.Awesome; v4 != nil {
		switch {
		case v4.stub != nil:
			return v4.invokeStub(m0)
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			return v4.invokeExpect(m0)
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m2.td.location, v3))
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

func Test_TargetData_implementationCode(t *testing.T) {
	cases := []struct {
		name       string
		receiver   string
		location   string
		skipExpect bool
		method     MethodInfo
		expected   string
	}{
		{
			name:     "empty",
			receiver: "m",
			location: "location",
			method: MethodInfo{
				Name:   "Empty",
				Struct: "targetEmpty",
			},
			expected: `package emitter

func (m *target) Empty() {
	v0, v1, v2 := "Target", "Empty", "()"
	v3 := []any{}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Empty; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub()
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect()
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name:     "only have results with no name originally",
			receiver: "m",
			location: "location",
			method: MethodInfo{
				Name:   "Method",
				Struct: "targetMethod",
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "", Type: gentest.Type("string")},
					{Name: "second", Field: "second", OriginalName: "", Type: gentest.Type("error")},
				},
			},
			expected: `package emitter

func (m *target) Method() (string, error) {
	v0, v1, v2 := "Target", "Method", "() (string, error)"
	v3 := []any{}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			return v4.invokeStub()
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			return v4.invokeExpect()
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name:     "has arguments only",
			receiver: "m",
			location: "location",
			method: MethodInfo{
				Name:   "Method",
				Struct: "targetMethod",
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
			},
			expected: `package emitter

import "context"

func (m *target) Method(ctx context.Context, input string) {
	v0, v1, v2 := "Target", "Method", "(ctx Context, input string)"
	v3 := []any{"ctx", ctx, "input", input}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			v4.invokeStub(ctx, input)
			return
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			v4.invokeExpect(ctx, input)
			return
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name:     "has arguments with named return",
			receiver: "m",
			location: "location",
			method: MethodInfo{
				Name:   "Method",
				Struct: "targetMethod",
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
				},
			},
			expected: `package emitter

import "context"

func (m *target) Method(ctx context.Context, input string) (err error) {
	v0, v1, v2 := "Target", "Method", "(ctx Context, input string) (err error)"
	v3 := []any{"ctx", ctx, "input", input}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			return v4.invokeStub(ctx, input)
		case len(v4.expects) > 0:
			v5 := len(v4.Calls)
			if v5 < len(v4.expects) {
				v4.expects[v5].tb.Helper()
			}
			return v4.invokeExpect(ctx, input)
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},

		{
			name:       "skip expect",
			receiver:   "m",
			location:   "location",
			skipExpect: true,
			method: MethodInfo{
				Name:   "Method",
				Struct: "targetMethod",
				Arguments: []VarInfo{
					{Name: "ctx", Field: "ctx", OriginalName: "ctx", Type: gentest.Type("context.Context")},
					{Name: "input", Field: "input", OriginalName: "input", Type: gentest.Type("string")},
				},
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "err", Type: gentest.Type("error")},
				},
			},
			expected: `package emitter

import "context"

func (m *target) Method(ctx context.Context, input string) (err error) {
	v0, v1, v2 := "Target", "Method", "(ctx Context, input string) (err error)"
	v3 := []any{"ctx", ctx, "input", input}

	if m.td == nil {
		panic(libMessageNotImplemented(v0, v1, v2, "", v3))
	}

	if v4 := m.td.Method; v4 != nil {
		switch {
		case v4.stub != nil:
			return v4.invokeStub(ctx, input)
		}
	}
	panic(libMessageNotImplemented(v0, v1, v2, m.td.location, v3))
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := TargetData{
				Interface:        "Target",
				Struct:           "target",
				Constructor:      "testTarget",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				ExpecterStruct:   "targetExpecter",
				Lib:              libData(),
				SkipExpect:       tc.skipExpect,
				Methods:          []MethodInfo{tc.method},
			}

			code := data.implementationCode(tc.receiver, tc.location, tc.method)
			jf := jen.NewFile("emitter")
			jf.Add(code)
			out := jf.GoString()
			assert.Equal(t, tc.expected, out)
		})
	}
}
