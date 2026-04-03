package mockgen

import "testing"

func Test_MethodExpecterValueData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodExpecterValueData
		expected string
	}{
		{
			name: "emit nothing if skip expect",
			data: MethodExpecterValueData{
				SkipExpect: true,
			},
			expected: "",
		},

		{
			name:     "emit nothing if there is no return",
			data:     MethodExpecterValueData{},
			expected: "",
		},

		{
			name: "emit a struct with return syntax enforced code",
			data: MethodExpecterValueData{
				ExpecterValueStruct: "targetMethodExpecterWithValue",
				ExpectStruct:        "targetMethodExpect",
				ReturnStruct:        "targetMethodReturn",
				Returns:             varInfos("First: first string", "Second: second error"),
			},
			expected: `package emitter

type targetMethodExpecterWithValue struct {
	expect *targetMethodExpect
}

func (e *targetMethodExpecterWithValue) Return(first string, second error) {
	e.expect.returns = targetMethodReturn{First: first, Second: second}
}
`,
		},

		{
			name: "emit a struct with receiver name aware about collision",
			data: MethodExpecterValueData{
				ExpecterValueStruct: "targetMethodExpecterWithValue",
				ExpectStruct:        "targetMethodExpect",
				ReturnStruct:        "targetMethodReturn",
				Returns:             varInfos("e: e string", "e0: e0 error", "e1: e1 error"),
			},
			expected: `package emitter

type targetMethodExpecterWithValue struct {
	expect *targetMethodExpect
}

func (e2 *targetMethodExpecterWithValue) Return(e string, e0 error, e1 error) {
	e2.expect.returns = targetMethodReturn{e: e, e0: e0, e1: e1}
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
