package mockgen

import (
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"
	genlib "nhatp.com/go/gen-lib"
)

func toPascalCase(name string) string {
	return strcase.ToGoPascal(name)
}

type VarInfo struct {
	Name         string
	Field        string
	OriginalName string
	Type         types.Type
}

type MethodInfo struct {
	Name                   string
	Struct                 string
	CallStruct             string
	ArgumentStruct         string
	ReturnStruct           string
	ArgumentMatcherStruct  string
	ExpectStruct           string
	ExpecterStruct         string
	ExpecterMatchStruct    string
	ExpecterMatchArgStruct string
	ExpecterValueStruct    string
	ExpecterValueArgStruct string
	Arguments              []VarInfo
	Returns                []VarInfo
}

func targetMethodSignatureString(method MethodInfo) string {
	signature := strings.Builder{}
	if len(method.Arguments) > 0 {
		signature.WriteRune('(')
		var strArgs []string
		for _, v := range method.Arguments {
			strArgs = append(strArgs, v.Name+" "+genlib.TypeSimpleName(v.Type))
		}
		signature.WriteString(strings.Join(strArgs, ", "))
		signature.WriteRune(')')
	} else {
		signature.WriteString("()")
	}

	if len(method.Returns) > 0 {
		wrap := false
		var strReturns []string
		for _, v := range method.Returns {
			if v.OriginalName != "" {
				wrap = true
				strReturns = append(strReturns, v.OriginalName+" "+genlib.TypeSimpleName(v.Type))
			} else {
				strReturns = append(strReturns, genlib.TypeSimpleName(v.Type))
			}
		}
		signature.WriteRune(' ')
		if wrap || len(strReturns) > 1 {
			signature.WriteRune('(')
			signature.WriteString(strings.Join(strReturns, ", "))
			signature.WriteRune(')')
		} else {
			signature.WriteString(strings.Join(strReturns, ","))
		}
	}
	return signature.String()
}

func targetMethodSignature(arguments []VarInfo, returns []VarInfo) jen.Code {
	var params, results []jen.Code

	for _, v := range arguments {
		if v.Name == "" {
			params = append(params, genlib.TypeToJenCode(v.Type))
		} else {
			params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		}
	}

	for _, v := range returns {
		if v.OriginalName == "" {
			results = append(results, genlib.TypeToJenCode(v.Type))
		} else {
			results = append(results, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
		}
	}

	return jen.Func().Params(params...).Params(results...)
}

func targetMethodMatcherSignature(args ...VarInfo) jen.Code {
	var params []jen.Code

	if len(args) == 0 {
		return nil
	}
	for _, v := range args {
		if v.Name == "" {
			params = append(params, genlib.TypeToJenCode(v.Type))
		} else {
			params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		}
	}
	return jen.Func().Params(params...).Params(jen.Bool())
}

// ---

type Emitter interface {
	Library(ctx genlib.EmitterContext, data LibraryData) []jen.Code
	Target(ctx EmitterContext, data TargetData) []jen.Code
	Stubber(ctx genlib.EmitterContext, data TargetStubberData) []jen.Code
	Expecter(ctx genlib.EmitterContext, data TargetExpecterData) []jen.Code
	Method(ctx genlib.EmitterContext, data MethodData) []jen.Code
	MethodExpecter(ctx genlib.EmitterContext, data MethodExpecterData) []jen.Code
	MethodExpecterMatch(ctx genlib.EmitterContext, data MethodExpecterMatchData) []jen.Code
	MethodExpecterMatchArg(ctx genlib.EmitterContext, data MethodExpecterMatchArgData) []jen.Code
	MethodExpecterValue(ctx genlib.EmitterContext, data MethodExpecterValueData) []jen.Code
	MethodExpecterValueArg(ctx genlib.EmitterContext, data MethodExpecterValueArgData) []jen.Code
}

type DefaultEmitter struct {
}

func (e *DefaultEmitter) Library(ctx genlib.EmitterContext, data LibraryData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) Target(ctx EmitterContext, data TargetData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) Stubber(ctx genlib.EmitterContext, data TargetStubberData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) Expecter(ctx genlib.EmitterContext, data TargetExpecterData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) Method(ctx genlib.EmitterContext, data MethodData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) MethodExpecter(ctx genlib.EmitterContext, data MethodExpecterData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) MethodExpecterMatch(ctx genlib.EmitterContext, data MethodExpecterMatchData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) MethodExpecterMatchArg(ctx genlib.EmitterContext, data MethodExpecterMatchArgData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) MethodExpecterValue(ctx genlib.EmitterContext, data MethodExpecterValueData) []jen.Code {
	return data.GenerateCode()
}

func (e *DefaultEmitter) MethodExpecterValueArg(ctx genlib.EmitterContext, data MethodExpecterValueArgData) []jen.Code {
	return data.GenerateCode()
}

var _ Emitter = (*DefaultEmitter)(nil)
