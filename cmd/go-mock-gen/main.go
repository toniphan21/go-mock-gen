package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	mockgen "nhatp.com/go/mock-gen"
)

type VersionCmd struct{}

type Arguments struct {
	Version *VersionCmd `arg:"subcommand:version" help:"Print version information and exit"`

	Interface   string `arg:"-i,--interface" placeholder:"NAME" help:"Comma-separated list of interfaces to mock (e.g. Repository,UserService)"`
	Struct      string `arg:"-s,--struct" placeholder:"STRUCT" help:"Struct name for the generated mock; only valid when mocking a single interface;\n                         defaults to the unexported interface name (e.g. Repository -> repository)"`
	PackageName string `arg:"-p,--package" placeholder:"PKG_NAME" help:"Package name for the generated code. Defaults to the source package name of the interface"`
	Output      string `arg:"-o,--output" placeholder:"PATH" help:"Output file for the generated code" default:"mockgen_test.go"`
	DryRun      bool   `arg:"-d,--dry-run" help:"Preview changes without writing to disk" default:"false"`
	OmitExpect  bool   `arg:"--omit-expect" help:"omit EXPECT mock generation" default:"false"`

	NoColor bool `arg:"--no-color" help:"Disable colors" default:"false"`
}

func (*Arguments) Epilogue() string {
	return `Examples:
  Generate mocks for a single interface:
    go-mock-gen -i Repository

  Generate mocks for multiple interfaces:
    go-mock-gen -i Repository,UserService

  Generate a mock with a custom struct name:
    go-mock-gen -i Repository -s repoMock

  Generate with a custom package and output file:
    go-mock-gen -i Repository -p mocks -o mocks/mockgen_test.go
`
}

func main() {
	var args Arguments
	p := arg.MustParse(&args)
	logger := slog.New(cli.NewSlogHandler(os.Stdout, "info"))

	switch {
	case args.Version != nil:
		v := fmt.Sprintf("%s%s - %s", mockgen.BinaryPath, color.Binary(mockgen.BinaryName), color.Version(mockgen.BinaryVersion))
		logger.Info(v)

	case args.Interface != "":
		generate(args, logger)

	default:
		p.WriteHelp(os.Stderr)
	}

	os.Exit(0)
}

func generate(cmd Arguments, logger *slog.Logger) {
	if cmd.NoColor {
		cli.DisableColor()
	}

	if cmd.DryRun {
		logger.Info(color.Binary(mockgen.BinaryName) + " " + color.Version(mockgen.BinaryVersion) + " in DRY mode")
	} else {
		logger.Info(color.Binary(mockgen.BinaryName) + " " + color.Version(mockgen.BinaryVersion))
	}

	var iface, ifaceStrings []string
	for _, v := range strings.Split(cmd.Interface, ",") {
		vv := strings.TrimSpace(v)
		if vv != "" {
			iface = append(iface, vv)
			ifaceStrings = append(ifaceStrings, color.Source(vv))
		}
	}

	if len(iface) == 0 {
		logger.Error(cli.ColorRed("no interface specified, use -i NAME (comma-separated list accepted)"))
	}
	logger.Info(fmt.Sprintf("generating mock for interface %s", strings.Join(ifaceStrings, ", ")))

	if len(iface) > 1 && strings.TrimSpace(cmd.Struct) != "" {
		logger.Error(cli.ColorRed("--struct/-s can only be used when generating a mock for a single interface"))
	}

	if len(iface) > 1 && strings.TrimSpace(cmd.PackageName) != "" {
		logger.Error(cli.ColorRed("--package/-p can only be used when generating a mock for a single interface"))
	}

	var structName, packageName, output string
	if strings.TrimSpace(cmd.Output) != "" {
		output = cmd.Output
	}

	if len(iface) == 1 {
		if strings.TrimSpace(cmd.Struct) != "" {
			structName = cmd.Struct
		}
		if strings.TrimSpace(cmd.PackageName) != "" {
			packageName = cmd.PackageName
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
			SkipExpect:    cmd.OmitExpect,
		})
	}

	dir, err := os.Getwd()
	if err != nil {
		logger.Error(cli.ColorRed(err.Error()))
		os.Exit(1)
	}

	fileManager := mockgen.NewFileManager(dir, mockgen.WithBinaryName(mockgen.BinaryFullName), mockgen.WithVersion(mockgen.BinaryVersion))
	generator := mockgen.New(fileManager)

	pkgs, err := mockgen.LoadPackages(dir)
	for _, pkg := range pkgs {
		var configs []mockgen.Config
		for _, v := range cfs {
			v.PackagePath = pkg.PkgPath
			configs = append(configs, v)
		}

		if err = generator.Generate(pkg, configs); err != nil {
			logger.Error(cli.ColorRed(err.Error()))
			os.Exit(1)
		}
	}

	if cmd.DryRun {
		logger.Info(color.Binary(mockgen.BinaryName) + " is printing generated file content")
		for _, out := range fileManager.Files() {
			content := out.Content()
			if cmd.NoColor {
				logger.Info(content)
			} else {
				cli.PrintFileWithFunction(out.RelPath, []byte(content), func(l string) {
					logger.Info(l)
				})
			}
		}
	} else {
		logger.Info(color.Binary(mockgen.BinaryName) + " is saving generated file to disk")
		for _, out := range fileManager.Files() {
			content := out.Content()
			outer := filepath.Dir(out.RelPath)
			if err = os.MkdirAll(outer, 0755); err != nil {
				logger.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}

			if err = os.WriteFile(out.FullPath, []byte(content), 0644); err != nil {
				logger.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}
			logger.Info(color.Binary(mockgen.BinaryName) + " saved " + out.RelPath)
		}
	}
	logger.Info(cli.ColorGreen("done"))
}
