package mockgen

import "testing"

func Test_TargetExpecterData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     TargetExpecterData
		expected string
	}{
		{
			name: "return nil if skip expect",
			data: TargetExpecterData{
				SkipExpect: true,
			},
			expected: ``,
		},

		{
			name:     "return nil if there is no method",
			data:     TargetExpecterData{},
			expected: ``,
		},

		{
			name: "strip the expected if it is skipped",
			data: TargetExpecterData{
				TargetStruct:                      "target",
				TargetTestDoubleStruct:            "targetTestDouble",
				TargetExpecterStruct:              "targetExpecter",
				TargetMethodStruct:                "targetMethod",
				TargetMethodExpectStruct:          "targetMethodExpect",
				TargetMethodExpecterStruct:        "targetMethodExpecter",
				TargetMethodArgumentMatcherStruct: "targetMethodArgumentMatcher",
				Lib:                               libData(),
				Methods: []MethodInfo{
					{
						Name:   "Method",
						Struct: "targetMethod",
					},
				},
			},
			expected: `package emitter

import "testing"

type targetExpecter struct {
	target *target
}

func (e *targetExpecter) Method(tb testing.TB) *targetMethodExpecter {
	if e.target.td == nil {
		e.target.td = &targetTestDouble{}
	}

	var m = e.target.td.Method
	if m == nil {
		m = &targetMethod{}
		e.target.td.Method = m
	}

	if m.stub != nil {
		m.panic(libMessageExpectAfterStub(m, m.stubLocation))
	}

	if tb == nil {
		m.panic(libMessageExpectByNil(m))
	}

	idx := len(m.expects)
	m.expects = append(m.expects, &targetMethodExpect{
		location:         libCallerLocation(2),
		matcher:          &targetMethodArgumentMatcher{},
		matcherWants:     make(map[string]any),
		matcherMethods:   make(map[string]string),
		matcherHints:     make(map[string]string),
		matcherLocations: make(map[string]string),
		index:            idx,
		tb:               tb,
	})

	tb.Helper()
	tb.Cleanup(func() {
		tb.Helper()
		m.verify(idx)
	})

	return &targetMethodExpecter{target: m, expect: m.expects[idx]}
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
