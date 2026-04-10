package yaml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mockgen "nhatp.com/go/mock-gen"
)

func errorMessage(lines ...string) string {
	return strings.Join(lines, "\n")
}

func ymlFile(lines ...string) string {
	return strings.Join(lines, "\n")
}

func TestParse_Invalid(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "invalid - unknown in global level",
			input: `unknown: 1`,
			expected: errorMessage(
				"yaml: unmarshal errors:",
				"  line 1: field unknown not found in type yaml.Global",
			),
		},

		{
			name:     "invalid - unknown version",
			input:    `version: 1.1`,
			expected: `unknown version, supported: "1"`,
		},

		{
			name:  "invalid - no package",
			input: `version: 1`,
			expected: errorMessage(
				"no packages",
				"",
				"  hint: add packages configuration",
				"    packages:",
				"      github.com/you/project/path:",
				"        - Repository # mocked interface",
			),
		},

		{
			name: "invalid - invalid package field",
			input: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - emit_example: Invalid",
			),
			expected: errorMessage(
				"yaml: unmarshal errors:",
				"  line 3: cannot unmarshal !!str `Invalid` into bool",
			),
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Parse([]byte(tc.input))
			require.NotNil(t, err)
			assert.Equal(t, tc.expected, err.Error())
			assert.Nil(t, out)
		})
	}
}

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected *Global
	}{
		{
			name: "empty",
			input: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
			),
			expected: &Global{
				Version: "1",
				All: All{
					Output: Output{File: "mockgen_test.go"},
				},
				Packages: map[string][]Package{
					"github.com/you/project/path": {
						{Interface: "Repository"},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Parse([]byte(tc.input))
			require.NoError(t, err)

			require.Equal(t, tc.expected, out)
		})
	}
}

func TestToConfig(t *testing.T) {
	cases := []struct {
		name     string
		yaml     string
		expected []mockgen.Config
	}{
		{
			name: "quick - 1 interface",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
				},
			},
		},

		{
			name: "detail - 1 interface - custom struct name",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - interface: Repository",
				"      struct: repository",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					StructName:    "repository",
				},
			},
		},

		{
			name: "detail - 1 interface - custom constructor name",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - interface: Repository",
				"      constructor: newRepository",
			),
			expected: []mockgen.Config{
				{
					PackagePath:     "github.com/you/project/path",
					Output:          mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName:   "Repository",
					ConstructorName: "newRepository",
				},
			},
		},

		{
			name: "detail - 1 interface - emit example",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - interface: Repository",
				"      emit_example: true",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					EmitExample:   true,
				},
			},
		},

		{
			name: "detail - 1 interface - omit expect",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - interface: Repository",
				"      omit_expect: true",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					OmitExpect:    true,
				},
			},
		},

		{
			name: "detail - 1 interface - override all output",
			yaml: ymlFile(
				"all:",
				"  output:",
				"    file: test.go",
				"",
				"packages:",
				"  github.com/you/project/path:",
				"    - interface: Repository",
				"      output:",
				"        pkg_name: mocktest",
				"        file: repo_test.go",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{PackageName: "mocktest", TestFileName: "repo_test.go"},
					InterfaceName: "Repository",
				},
			},
		},

		{
			name: "all - use global emit example",
			yaml: ymlFile(
				"all:",
				"  emit_example: true",
				"",
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					EmitExample:   true,
				},
			},
		},

		{
			name: "all - use global emit example - entry override",
			yaml: ymlFile(
				"all:",
				"  emit_example: true",
				"",
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
				"    - interface: Service",
				"      emit_example: false",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					EmitExample:   true,
				},
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Service",
				},
			},
		},

		{
			name: "all - use global omit expect",
			yaml: ymlFile(
				"all:",
				"  omit_expect: true",
				"",
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					OmitExpect:    true,
				},
			},
		},

		{
			name: "all - use global omit expect - entry override",
			yaml: ymlFile(
				"all:",
				"  omit_expect: true",
				"",
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
				"    - interface: Service",
				"      omit_expect: false",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
					OmitExpect:    true,
				},
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Service",
				},
			},
		},

		{
			name: "quick - 2 interfaces",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/path:",
				"    - Repository",
				"    - Service",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
				},
				{
					PackagePath:   "github.com/you/project/path",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Service",
				},
			},
		},

		{
			name: "quick - 2 packages",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/infra:",
				"    - Repository",
				"  github.com/you/project/domain:",
				"    - Service",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/domain",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Service",
				},
				{
					PackagePath:   "github.com/you/project/infra",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Repository",
				},
			},
		},

		{
			name: "mixed - 2 packages",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/infra:",
				"    - interface: Repository",
				"      struct: repository",
				"      constructor: newRepository",
				"  github.com/you/project/domain:",
				"    - Service",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/domain",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "Service",
				},
				{
					PackagePath:     "github.com/you/project/infra",
					Output:          mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName:   "Repository",
					StructName:      "repository",
					ConstructorName: "newRepository",
				},
			},
		},

		{
			name: "mixed - multiple packages",
			yaml: ymlFile(
				"packages:",
				"  github.com/you/project/a:",
				"    - A",
				"  github.com/you/project/b:",
				"    - B",
				"  github.com/you/project/c:",
				"    - C",
				"  github.com/you/project/d:",
				"    - D",
				"  github.com/you/project/e:",
				"    - E",
				"  github.com/you/project/f:",
				"    - F",
			),
			expected: []mockgen.Config{
				{
					PackagePath:   "github.com/you/project/a",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "A",
				},
				{
					PackagePath:   "github.com/you/project/b",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "B",
				},
				{
					PackagePath:   "github.com/you/project/c",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "C",
				},
				{
					PackagePath:   "github.com/you/project/d",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "D",
				},
				{
					PackagePath:   "github.com/you/project/e",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "E",
				},
				{
					PackagePath:   "github.com/you/project/f",
					Output:        mockgen.Output{TestFileName: "mockgen_test.go"},
					InterfaceName: "F",
				},
			},
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := Parse([]byte(tc.yaml))
			require.NoError(t, err)
			assert.Equal(t, tc.expected, g.ToConfig())
		})
	}
}
