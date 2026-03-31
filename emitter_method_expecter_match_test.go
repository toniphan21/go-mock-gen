package mockgen

import "testing"

func Test_MethodExpecterMatchData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     MethodExpecterMatchData
		expected string
	}{
		{
			name: "emit nothing if skip expect",
			data: MethodExpecterMatchData{
				SkipExpect: true,
			},
			expected: "",
		},

		{
			name:     "emit nothing if there is no return",
			data:     MethodExpecterMatchData{},
			expected: "",
		},

		{
			name: "emit a struct with return syntax enforced code",
			data: MethodExpecterMatchData{
				TargetMethodExpecterMatchStruct: "targetMethodExpecterWithMatch",
				TargetMethodExpectStruct:        "targetMethodExpect",
				TargetMethodReturnStruct:        "targetMethodReturn",
				Returns:                         varInfos("First: first string", "Second: second error"),
			},
			expected: `package emitter

type targetMethodExpecterWithMatch struct {
	expect *targetMethodExpect
}

func (e *targetMethodExpecterWithMatch) Return(first string, second error) {
	e.expect.returns = targetMethodReturn{First: first, Second: second}
}
`,
		},

		{
			name: "emit a struct with receiver name aware about collision",
			data: MethodExpecterMatchData{
				TargetMethodExpecterMatchStruct: "targetMethodExpecterWithMatch",
				TargetMethodExpectStruct:        "targetMethodExpect",
				TargetMethodReturnStruct:        "targetMethodReturn",
				Returns:                         varInfos("e: e string", "e0: e0 error", "e1: e1 error"),
			},
			expected: `package emitter

type targetMethodExpecterWithMatch struct {
	expect *targetMethodExpect
}

func (e2 *targetMethodExpecterWithMatch) Return(e string, e0 error, e1 error) {
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
