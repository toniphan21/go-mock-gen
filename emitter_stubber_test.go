package mockgen

import (
	"testing"
)

func Test_TargetStubberData_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     TargetStubberData
		expected string
	}{
		{
			name:     "return nil if there is no method",
			data:     TargetStubberData{},
			expected: ``,
		},

		{
			name: "strip the expected if it is skipped",
			data: TargetStubberData{
				Struct:           "target",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:   "Method",
						Struct: "targetMethod",
					},
				},
				SkipExpect: true,
			},
			expected: `package emitter

type targetStubber struct {
	target *target
}

func (s *targetStubber) Method(stub func()) *targetMethod {
	if s.target.td == nil {
		s.target.td = &targetTestDouble{}
	}

	m := s.target.td.Method
	if m == nil {
		m = &targetMethod{stubLocation: libCallerLocation(2)}
		s.target.td.Method = m
	}

	if stub == nil {
		m.panic(libMessageStubByNil(m, libCallerLocation(2)))
	}

	if m.stub != nil {
		m.panic(libMessageDuplicateStub(m, m.stubLocation))
	}

	m.stub = stub
	return m
}
`,
		},

		{
			name: "should handle correct signature",
			data: TargetStubberData{
				Struct:           "target",
				TestDoubleStruct: "targetTestDouble",
				StubberStruct:    "targetStubber",
				Lib:              libData(),
				Methods: []MethodInfo{
					{
						Name:      "Method",
						Struct:    "targetMethod",
						Arguments: varInfos("Input: input string"),
						Returns:   varInfos("Output: output string", "Error: err error"),
					},
				},
			},
			expected: `package emitter

type targetStubber struct {
	target *target
}

func (s *targetStubber) Method(stub func(input string) (string, error)) *targetMethod {
	if s.target.td == nil {
		s.target.td = &targetTestDouble{}
	}

	m := s.target.td.Method
	if m == nil {
		m = &targetMethod{stubLocation: libCallerLocation(2)}
		s.target.td.Method = m
	}

	if stub == nil {
		m.panic(libMessageStubByNil(m, libCallerLocation(2)))
	}

	if m.stub != nil {
		m.panic(libMessageDuplicateStub(m, m.stubLocation))
	}

	if len(m.expects) > 0 {
		m.panic(libMessageStubAfterExpect(m, m.expects[0].location))
	}

	m.stub = stub
	return m
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
