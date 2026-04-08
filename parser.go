package mockgen

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
)

type parsedMethod struct {
	name      string
	arguments []*types.Var
	results   []*types.Var
}

type interfaceInfo struct {
	pkgPath string
	name    string
}

func getInterfaceInfo(input string) interfaceInfo {
	result := interfaceInfo{}

	s := input

	lastSlash := strings.LastIndex(s, "/")
	separatorIndex := strings.LastIndex(s, ".")

	if separatorIndex > lastSlash {
		result.pkgPath = s[:separatorIndex]
		s = s[separatorIndex+1:]
	} else {
		result.pkgPath = ""
	}

	result.name = s
	return result
}

func parse(dir string, pkg *packages.Package, iface interfaceInfo, namer Namer) []MethodInfo {
	var methods []parsedMethod
	if iface.pkgPath == "" {
		methods = parseSignatures(pkg, iface.name)
	} else {
		cf := genlib.LoadPackagesConfig(dir)
		pkgs, err := packages.Load(cf, iface.pkgPath)
		if err != nil || len(pkgs) != 1 {
			return nil
		}
		methods = parseSignatures(pkgs[0], iface.name)
	}

	if methods == nil {
		return nil
	}

	var out []MethodInfo
	for _, method := range methods {
		m := namer.Method(method.name)
		info := &MethodInfo{
			Name:                   method.name,
			Struct:                 m.Struct(),
			CallStruct:             m.Call(),
			ArgumentStruct:         m.Argument(),
			ArgumentMatcherStruct:  m.ArgumentMatcher(),
			ReturnStruct:           m.Return(),
			ExpectStruct:           m.Expect(),
			ExpecterStruct:         m.Expecter(),
			ExpecterMatchStruct:    m.ExpecterMatch(),
			ExpecterMatchArgStruct: m.ExpecterMatchArg(),
			ExpecterValueStruct:    m.ExpecterValue(),
			ExpecterValueArgStruct: m.ExpecterValueArg(),
		}

		for i, arg := range method.arguments {
			field, name := m.ArgumentField(arg.Name(), i)
			info.Arguments = append(info.Arguments, VarInfo{
				Field:        field,
				Name:         name,
				OriginalName: arg.Name(),
				Type:         arg.Type(),
			})
		}

		for i, ret := range method.results {
			field, name := m.ReturnField(ret.Name(), i)
			info.Returns = append(info.Returns, VarInfo{
				Field:        field,
				Name:         name,
				OriginalName: ret.Name(),
				Type:         ret.Type(),
			})
		}
		out = append(out, *info)
	}
	return out
}

func parseSignatures(pkg *packages.Package, interfaceName string) []parsedMethod {
	obj := pkg.Types.Scope().Lookup(interfaceName)
	if obj == nil {
		return nil
	}

	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil
	}

	var methods []parsedMethod
	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		sig, ok := m.Type().(*types.Signature)
		if !ok {
			continue
		}

		currentMethod := parsedMethod{
			name:      m.Name(),
			arguments: extractVars(sig.Params()),
			results:   extractVars(sig.Results()),
		}
		methods = append(methods, currentMethod)
	}
	return methods
}

func extractVars(tuple *types.Tuple) []*types.Var {
	if tuple == nil {
		return nil
	}
	vars := make([]*types.Var, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		vars[i] = tuple.At(i)
	}
	return vars
}
