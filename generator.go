package mockgen

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
)

type generatorImpl struct {
	fileManager     FileManager
	emitter         Emitter
	pkgNameManagers map[string]NameManager
	logger          Logger
}

func filterConfigs(pkg *packages.Package, configs []Config) []Config {
	var out []Config
	for _, v := range configs {
		if v.PackagePath != pkg.PkgPath {
			continue
		}

		if v.InterfaceName == "" {
			continue
		}

		c := v
		if c.Namer == nil {
			c.Namer = NewNamer(
				v.InterfaceName,
				WithStructName(v.StructName),
			)
		}
		out = append(out, c)
	}
	return out
}

func (g *generatorImpl) Generate(pkg *packages.Package, configs []Config) error {
	matchPkgConfigs := filterConfigs(pkg, configs)

	cfs := make(map[Config][]MethodInfo)
	for _, config := range matchPkgConfigs {
		methods := parse(pkg, config.InterfaceName, config.Namer)
		if len(methods) == 0 {
			continue
		}
		cfs[config] = methods
	}
	g.logger.Parsed(pkg, cfs)

	for config, methods := range cfs {
		if err := g.generate(pkg, config, methods); err != nil {
			g.logger.Error("generate", err)
			return err
		}
	}
	return nil
}

func (g *generatorImpl) getNameManager(sourcePkg *packages.Package, gf *genlib.GenFile) (NameManager, bool) {
	if gf.PkgPath == sourcePkg.PkgPath {
		nm, have := g.pkgNameManagers[gf.PkgPath]
		if !have {
			nm = genlib.NewNameManager("typ", sourcePkg.Types.Scope().Names())
			g.pkgNameManagers[gf.PkgPath] = nm
		}
		return nm, have
	}

	nm, have := g.pkgNameManagers[gf.PkgPath]
	if have {
		return nm, have
	}

	pkgs, err := packages.Load(genlib.LoadPackagesConfig(g.fileManager.RootDir()), gf.PkgPath)
	if err != nil || len(pkgs) != 1 {
		nm = genlib.NewNameManager("typ", nil)
	} else {
		nm = genlib.NewNameManager("typ", pkgs[0].Types.Scope().Names())
	}
	g.pkgNameManagers[gf.PkgPath] = nm
	return nm, have
}

func (g *generatorImpl) generate(pkg *packages.Package, config Config, info []MethodInfo) error {
	gf, err := g.fileManager.TestFile(pkg, config.Output)
	if err != nil {
		return err
	}
	ctx := NewEmitterContext(pkg, g.fileManager, gf, "v")
	g.logger.StartGenerating(pkg, config)

	var lib = config.Namer.Library()
	nm, have := g.getNameManager(pkg, gf)
	if !have {
		libNames := config.Namer.Library()
		lib = LibraryData{
			CallerLocationFunc:            nm.Request(libNames.CallerLocationFunc),
			MethodInterface:               nm.Request(libNames.MethodInterface),
			MessageWriteArgumentsFunc:     nm.Request(libNames.MessageWriteArgumentsFunc),
			MessageMatchFailFunc:          nm.Request(libNames.MessageMatchFailFunc),
			MessageNotImplementedFunc:     nm.Request(libNames.MessageNotImplementedFunc),
			MessageCallHistoryFunc:        nm.Request(libNames.MessageCallHistoryFunc),
			MessageTooManyCallsFunc:       nm.Request(libNames.MessageTooManyCallsFunc),
			MessageMatchByNilFunc:         nm.Request(libNames.MessageMatchByNilFunc),
			MessageExpectByNilFunc:        nm.Request(libNames.MessageExpectByNilFunc),
			MessageExpectAfterStubFunc:    nm.Request(libNames.MessageExpectAfterStubFunc),
			MessageStubByNilFunc:          nm.Request(libNames.MessageStubByNilFunc),
			MessageStubAfterExpectFunc:    nm.Request(libNames.MessageStubAfterExpectFunc),
			MessageDuplicateStubFunc:      nm.Request(libNames.MessageDuplicateStubFunc),
			MessageExpectButNotCalledFunc: nm.Request(libNames.MessageExpectButNotCalledFunc),
			MessageMatchArgByNilFunc:      nm.Request(libNames.MessageMatchArgByNilFunc),
			MessageDuplicateMatchArgFunc:  nm.Request(libNames.MessageDuplicateMatchArgFunc),
			MessageMatchArgHintFunc:       nm.Request(libNames.MessageMatchArgHintFunc),
			MatchArgumentFunc:             nm.Request(libNames.MatchArgumentFunc),
			ReflectEqualMatcherFunc:       nm.Request(libNames.ReflectEqualMatcherFunc),
			BasicComparisonMatcherFunc:    nm.Request(libNames.BasicComparisonMatcherFunc),
		}

		g.collect(gf, g.emitter.Library(ctx, lib))
	}

	targetStruct := nm.Request(config.Namer.Struct())
	targetConstructor := nm.Request(config.Namer.Constructor())
	targetTestDoubleStruct := nm.Request(config.Namer.TestDouble())
	targetStubberStruct := nm.Request(config.Namer.Stubber())
	targetExpecterStruct := nm.Request(config.Namer.Expecter())

	var methods []MethodInfo
	for _, v := range info {
		reg := v
		reg.Struct = nm.Request(v.Struct)
		reg.CallStruct = nm.Request(v.CallStruct)
		reg.ArgumentStruct = nm.Request(v.ArgumentStruct)
		reg.ReturnStruct = nm.Request(v.ReturnStruct)
		reg.ArgumentMatcherStruct = nm.Request(v.ArgumentMatcherStruct)
		reg.ExpectStruct = nm.Request(v.ExpectStruct)
		reg.ExpecterStruct = nm.Request(v.ExpecterStruct)
		reg.ExpecterMatchStruct = nm.Request(v.ExpecterMatchStruct)
		reg.ExpecterMatchArgStruct = nm.Request(v.ExpecterMatchArgStruct)
		reg.ExpecterValueStruct = nm.Request(v.ExpecterValueStruct)
		reg.ExpecterValueArgStruct = nm.Request(v.ExpecterValueArgStruct)
		methods = append(methods, reg)
	}

	g.collect(gf, g.emitter.Target(ctx, TargetData{
		Interface:        config.InterfaceName,
		Struct:           targetStruct,
		Constructor:      targetConstructor,
		TestDoubleStruct: targetTestDoubleStruct,
		StubberStruct:    targetStubberStruct,
		ExpecterStruct:   targetExpecterStruct,
		Lib:              lib,
		Methods:          methods,
		SkipExpect:       config.SkipExpect,
	}))

	g.collect(gf, g.emitter.Stubber(ctx, TargetStubberData{
		Struct:           targetStruct,
		StubberStruct:    targetStubberStruct,
		TestDoubleStruct: targetTestDoubleStruct,
		Methods:          methods,
		Lib:              lib,
		SkipExpect:       config.SkipExpect,
	}))

	g.collect(gf, g.emitter.Expecter(ctx, TargetExpecterData{
		Struct:           targetStruct,
		ExpecterStruct:   targetExpecterStruct,
		TestDoubleStruct: targetTestDoubleStruct,
		Methods:          methods,
		Lib:              lib,
		SkipExpect:       config.SkipExpect,
	}))

	for _, method := range methods {
		g.collect(gf, g.emitter.Method(ctx, MethodData{
			Struct:                method.Struct,
			CallStruct:            method.CallStruct,
			ArgumentStruct:        method.ArgumentStruct,
			ArgumentMatcherStruct: method.ArgumentMatcherStruct,
			ReturnStruct:          method.ReturnStruct,
			ExpectStruct:          method.ExpectStruct,
			Interface:             config.InterfaceName,
			Name:                  method.Name,
			Arguments:             method.Arguments,
			Returns:               method.Returns,
			Lib:                   lib,
			SkipExpect:            config.SkipExpect,
		}))

		g.collect(gf, g.emitter.MethodExpecter(ctx, MethodExpecterData{
			ExpectStruct:           method.ExpectStruct,
			ExpecterStruct:         method.ExpecterStruct,
			ExpecterMatchStruct:    method.ExpecterMatchStruct,
			ExpecterMatchArgStruct: method.ExpecterMatchArgStruct,
			ExpecterValueStruct:    method.ExpecterValueStruct,
			ExpecterValueArgStruct: method.ExpecterValueArgStruct,
			Struct:                 method.Struct,
			ReturnStruct:           method.ReturnStruct,
			Arguments:              method.Arguments,
			Returns:                method.Returns,
			Lib:                    lib,
			SkipExpect:             config.SkipExpect,
		}))

		g.collect(gf, g.emitter.MethodExpecterMatch(ctx, MethodExpecterMatchData{
			ExpecterMatchStruct: method.ExpecterMatchStruct,
			ExpectStruct:        method.ExpectStruct,
			ReturnStruct:        method.ReturnStruct,
			Returns:             method.Returns,
			SkipExpect:          config.SkipExpect,
		}))

		g.collect(gf, g.emitter.MethodExpecterMatchArg(ctx, MethodExpecterMatchArgData{
			ExpecterMatchArgStruct: method.ExpecterMatchArgStruct,
			ExpectStruct:           method.ExpectStruct,
			Struct:                 method.Struct,
			ReturnStruct:           method.ReturnStruct,
			Arguments:              method.Arguments,
			Returns:                method.Returns,
			Lib:                    lib,
			SkipExpect:             config.SkipExpect,
		}))

		g.collect(gf, g.emitter.MethodExpecterValue(ctx, MethodExpecterValueData{
			ExpecterValueStruct: method.ExpecterValueStruct,
			ExpectStruct:        method.ExpectStruct,
			ReturnStruct:        method.ReturnStruct,
			Returns:             method.Returns,
			SkipExpect:          config.SkipExpect,
		}))

		g.collect(gf, g.emitter.MethodExpecterValueArg(ctx, MethodExpecterValueArgData{
			ExpecterValueArgStruct: method.ExpecterValueArgStruct,
			ExpectStruct:           method.ExpectStruct,
			Struct:                 method.Struct,
			ReturnStruct:           method.ReturnStruct,
			Arguments:              method.Arguments,
			Returns:                method.Returns,
			Lib:                    lib,
			SkipExpect:             config.SkipExpect,
		}))

		if config.EmitExamples {
			eo := g.makeExampleOutput(config.Output)
			egf, err := g.fileManager.Make(pkg, eo.PackageName, eo.TestFileName)
			if err != nil {
				return err
			}

			g.collect(egf, g.emitter.Example(ctx, ExampleData{
				Constructor:   targetConstructor,
				InterfaceName: config.InterfaceName,
				MethodName:    method.Name,
				Arguments:     method.Arguments,
				Returns:       method.Returns,
				SkipExpect:    config.SkipExpect,
			}))
		}
	}

	g.logger.Generated(pkg, GeneratedInfo{
		Config:      config,
		Struct:      targetStruct,
		Constructor: targetConstructor,
	})
	return nil
}

func (g *generatorImpl) collect(gf *genlib.GenFile, codes []jen.Code) {
	jf := gf.JenFile()
	for _, code := range codes {
		if code != nil {
			jf.Add(code)
		}
	}
}

func (g *generatorImpl) makeExampleOutput(output Output) Output {
	var testFileName = output.TestFileName
	switch {
	case strings.HasSuffix(testFileName, "_test.go"):
		testFileName = strings.TrimSuffix(testFileName, "_test.go") + "_example_test.go"

	case strings.HasSuffix(testFileName, ".go"):
		testFileName = strings.TrimSuffix(testFileName, ".go") + "_example.go"

	default:
		testFileName = testFileName + ".example"
	}

	return Output{
		PackageName:  output.PackageName,
		TestFileName: testFileName,
	}
}
