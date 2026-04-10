package cmd

import (
	"errors"
	"strings"

	"nhatp.com/go/gen-lib/cli/color"
	mockgen "nhatp.com/go/mock-gen"
)

type VersionCmd struct{}

type TestCmd struct {
	Files     []string `arg:"positional" help:"markdown file(s) to test" placeholder:"FILE"`
	Name      string   `arg:"-n,--name" help:"run test which has matched name (case insensitive)" default:""`
	ShowSetup bool     `arg:"-u,--show-setup" help:"show test setup steps" default:"false"`
	TabSize   int      `arg:"-t,--tab-size" help:"number of spaces to use in tab size" default:"8"`
	EmitCode  string   `arg:"-e,--emit-code" help:"emit to code if the test passed. If empty looking for path in Markdown comment." default:""`
}

type Arguments struct {
	Version *VersionCmd `arg:"subcommand:version" help:"print version information and exit"`
	Test    *TestCmd    `arg:"subcommand:test" help:"test generator using markdown files"`

	Interface   string `arg:"-i,--interface" placeholder:"NAMES" help:"comma-separated list of interfaces to mock. Supports:\n                           - local    : Repository\n                           - qualified: io.Reader\n                           - full path: github.com/user/pkg.Interface"`
	Struct      string `arg:"-s,--struct" placeholder:"STRUCT" help:"struct name for the generated mock; only valid when mocking a single interface;\n                         defaults to the unexported interface name (e.g. Repository -> repository)"`
	PackageName string `arg:"-p,--package" placeholder:"PKG_NAME" help:"package name for the generated code. Defaults to the source package name of the interface"`
	Output      string `arg:"-o,--output" placeholder:"PATH" help:"output file for the generated code" default:"mockgen_test.go"`
	DryRun      bool   `arg:"-d,--dry-run" help:"preview changes without writing to disk" default:"false"`
	EmitExample bool   `arg:"--example" help:"emit test examples" default:"false"`
	OmitExpect  bool   `arg:"--omit-expect" help:"omit EXPECT mock generation" default:"false"`

	NoColor bool `arg:"--no-color" help:"disable colors" default:"false"`
	Verbose bool `arg:"-v,--verbose" help:"enable verbose logging"`
}

func (*Arguments) Epilogue() string {
	return `Examples:
  # Generate a mock for a local interface:
  go-mock-gen -i Repository

  # Generate a mock for a single interface with example tests:
  go-mock-gen -i Repository --example

  # Generate a mock from standard library or external packages:
  go-mock-gen -i io.Reader,net/http.RoundTripper

  # Generate mocks for multiple local interfaces:
  go-mock-gen -i Repository,UserService

  # Generate a mock with a custom struct name:
  go-mock-gen -i Repository -s repoMock

  # Generate with a custom package and output file:
  go-mock-gen -i Repository -s Repository -p mock -o mock/mockgen_test.go
`
}

func ToConfigs(dir string, args Arguments) ([]mockgen.Config, error) {
	var iface, ifaceStrings []string
	for _, v := range strings.Split(args.Interface, ",") {
		vv := strings.TrimSpace(v)
		if vv != "" {
			iface = append(iface, vv)
			ifaceStrings = append(ifaceStrings, color.Source(vv))
		}
	}

	if len(iface) == 0 {
		return nil, errors.New("no interface specified, use -i NAME (comma-separated list accepted)")
	}

	if len(iface) > 1 && strings.TrimSpace(args.Struct) != "" {
		return nil, errors.New("--struct/-s can only be used when generating a mock for a single interface")
	}

	if len(iface) > 1 && strings.TrimSpace(args.PackageName) != "" {
		return nil, errors.New("--package/-p can only be used when generating a mock for a single interface")
	}

	var structName, packageName, output string
	if strings.TrimSpace(args.Output) != "" {
		output = args.Output
	}

	if len(iface) == 1 {
		if strings.TrimSpace(args.Struct) != "" {
			structName = args.Struct
		}
		if strings.TrimSpace(args.PackageName) != "" {
			packageName = args.PackageName
		}
	}

	var cfs []mockgen.Config
	for _, v := range iface {
		cfs = append(cfs, mockgen.Config{
			Output: mockgen.Output{
				PackageName:  packageName,
				TestFileName: output,
			},
			InterfaceName: v,
			StructName:    structName,
			SkipExpect:    args.OmitExpect,
			EmitExamples:  args.EmitExample,
		})
	}

	var configs []mockgen.Config
	pkgs, err := mockgen.LoadPackages(dir)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		if pkg.Dir != dir {
			continue
		}
		for _, v := range cfs {
			v.PackagePath = pkg.PkgPath
			configs = append(configs, v)
		}
	}
	return configs, nil
}
