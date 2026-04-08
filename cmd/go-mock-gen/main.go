package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	mockgen "nhatp.com/go/mock-gen"
	"nhatp.com/go/mock-gen/internal/cmd"
	"nhatp.com/go/mock-gen/mockgentest"
)

func main() {
	var args cmd.Arguments
	p := arg.MustParse(&args)

	if args.NoColor {
		cli.DisableColor()
	}

	level := "info"
	if args.Verbose {
		level = "debug"
	}
	logger := slog.New(cli.NewSlogHandler(os.Stdout, level))

	switch {
	case args.Version != nil:
		v := fmt.Sprintf("%s%s - %s", mockgen.BinaryPath, color.Binary(mockgen.BinaryName), color.Version(mockgen.BinaryVersion))
		logger.Info(v)

	case args.Test != nil:
		test(&cli.TestRunner{
			Files:     args.Test.Files,
			Name:      args.Test.Name,
			TabSize:   args.Test.TabSize,
			ShowSetup: args.Test.ShowSetup,
			EmitPath:  args.Test.EmitCode,
			Logger:    logger,
			FilePathResolver: cli.WithVanityURLFilePathResolver(map[string]string{
				"repo://":      mockgen.RawGitRefsURL + "/heads/main/",
				"repo-refs://": mockgen.RawGitRefsURL,
			}),
		})

	default:
		if canGenerate(args) {
			generate(args, logger)
		} else {
			p.WriteHelp(os.Stderr)
		}
	}

	os.Exit(0)
}

func canGenerate(args cmd.Arguments) bool {
	return args.Interface != ""
}

func generate(args cmd.Arguments, logger *slog.Logger) {
	if args.DryRun {
		logger.Info(color.Binary(mockgen.BinaryName) + " " + color.Version(mockgen.BinaryVersion) + " in DRY mode")
	} else {
		logger.Info(color.Binary(mockgen.BinaryName) + " " + color.Version(mockgen.BinaryVersion))
	}

	dir, err := os.Getwd()
	if err != nil {
		logger.Error(cli.ColorRed(err.Error()))
		os.Exit(1)
	}

	logPoints := &mockgen.LogPoints{
		Error: func(action string, err error) {
			logger.Error(cli.ColorRed(action + ": " + err.Error()))
		},
		FilteredConfigs: func(pkg *packages.Package, parsed map[mockgen.Config][]mockgen.MethodInfo) {
			var interfaces []string
			for config := range parsed {
				interfaces = append(interfaces, color.Source(config.InterfaceName))
			}
			if len(interfaces) > 0 {
				logger.Info(fmt.Sprintf("package %s:", color.Package(pkg.PkgPath)))
			}
		},
		Generated: func(pkg *packages.Package, info mockgen.GeneratedInfo) {
			logger.Info(fmt.Sprintf(
				"\t generated mock struct %s for interface %s with constructor %s",
				color.Generated(info.Struct),
				color.Source(info.Config.InterfaceName),
				color.Generated(info.Constructor+"()"),
			))
		},
	}

	fileManager := mockgen.NewFileManager(dir, mockgen.WithBinaryName(mockgen.BinaryFullName), mockgen.WithVersion(mockgen.BinaryVersion))
	generator := mockgen.New(fileManager, mockgen.WithLogPoints(logPoints))

	pkgs, err := mockgen.LoadPackages(dir)
	if err != nil {
		logger.Error(cli.ColorRed(err.Error()))
		os.Exit(1)
	}
	configs, err := cmd.ToConfigs(dir, args)
	if err != nil {
		logger.Error(cli.ColorRed(err.Error()))
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		if err = generator.Generate(pkg, configs); err != nil {
			os.Exit(1)
		}
	}

	if args.DryRun {
		logger.Info(color.Binary(mockgen.BinaryName) + " is printing generated file content")
		for _, out := range fileManager.Files() {
			content := out.Content()
			if args.NoColor {
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
			logger.Info(color.Binary(mockgen.BinaryName) + " saved " + color.Generated(out.RelPath))
		}
	}
	logger.Info(cli.ColorGreen("done"))
}

func test(cmd *cli.TestRunner) {
	cmd.Print("Running tests with " + color.Binary(mockgen.BinaryName) + " " + color.Version(mockgen.BinaryVersion))
	cmd.Print("")

	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		dir := tc.TestDir
		gtc := mockgentest.GoldenTestCase{
			Source:      tc.SourceFiles,
			GoldenFiles: tc.GoldenFiles,
		}

		err := gtc.Setup(dir)
		if err != nil {
			cmd.PrintWarn("\tcannot setup " + err.Error())
			return nil, err
		}

		configs, err := gtc.ParseConfigs(dir)
		if err != nil {
			cmd.PrintWarn("\tcannot parse configs " + err.Error())
			return nil, err
		}

		var opts []mockgen.Option
		return gtc.ExecuteMockGen(dir, configs, opts...)
	}
	cmd.Run()
}
