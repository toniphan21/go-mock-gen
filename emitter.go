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
	Name      string
	Struct    string
	Arguments []VarInfo
	Returns   []VarInfo
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
