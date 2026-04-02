package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

func targetMethodExpecterReturnCode(receiverName, receiverType, targetMethodReturnStruct string, returns []VarInfo) jen.Code {
	var params []jen.Code
	var values []jen.Code
	for _, v := range returns {
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		values = append(values, jen.Id(v.Field).Op(":").Id(v.Name))
	}

	return jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id("Return").
		Params(params...).
		Block(
			jen.Id(receiverName).Dot("expect").Dot("returns").Op("=").Id(targetMethodReturnStruct).Values(values...),
		).
		Line()
}

type MethodExpecterData struct {
	TargetMethodExpectStruct           string
	TargetMethodExpecterStruct         string
	TargetMethodExpecterMatchStruct    string
	TargetMethodExpecterMatchArgStruct string
	TargetMethodExpecterValueStruct    string
	TargetMethodExpecterValueArgStruct string
	TargetMethodStruct                 string
	TargetMethodReturnStruct           string
	Arguments                          []VarInfo
	Returns                            []VarInfo
	Lib                                LibraryData
	SkipExpect                         bool
}

func (d *MethodExpecterData) structCode() jen.Code {
	fields := []jen.Code{
		jen.Id("expect").Op("*").Id(d.TargetMethodExpectStruct),
	}
	if len(d.Arguments) > 0 {
		fields = append(fields, jen.Id("target").Op("*").Id(d.TargetMethodStruct))
	}
	return jen.Type().Id(d.TargetMethodExpecterStruct).Struct(fields...).Line()
}

func (d *MethodExpecterData) matchFuncCode(receiver string) jen.Code {
	var signature = targetMethodMatcherSignature(d.Arguments...)

	return jen.Func().
		Params(jen.Id(receiver).Op("*").Id(d.TargetMethodExpecterStruct)).
		Id("Match").Params(jen.Id("matcher").Add(signature)).
		Op("*").Id(d.TargetMethodExpecterMatchStruct).
		Block(
			jen.If(jen.Id("matcher").Op("==").Nil()).Block(
				jen.Id(receiver).Dot("expect").Dot("tb").Dot("Helper").Call(),
				jen.Id(receiver).Dot("target").Dot("fatal").Call(
					jen.Id(receiver).Dot("expect").Dot("index"),
					jen.Id(d.Lib.MessageMatchByNilFunc).Call(jen.Id(receiver).Dot("target")),
				),
			),
			jen.Line(),
			jen.Id(receiver).Dot("expect").Dot("match").Op("=").Id("matcher"),
			jen.Id(receiver).Dot("expect").Dot("matchLocation").Op("=").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2)),
			jen.Return(
				jen.Op("&").Id(d.TargetMethodExpecterMatchStruct).Values(
					jen.Id("expect").Op(":").Id(receiver).Dot("expect"),
				),
			),
		).Line()
}

func (d *MethodExpecterData) withFuncCode(receiver string) jen.Code {
	var params, body []jen.Code
	for _, v := range d.Arguments {
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))

		fn := "With" + toPascalCase(v.Name)

		body = append(body, jen.Id(receiver).Dot(fn).Call(jen.Id(v.Name)))
		body = append(body,
			jen.Id(receiver).Dot("expect").Dot("matcherLocations").
				Index(jen.Lit(v.Name)).Op("=").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2)),
		)
		body = append(body, jen.Line())
	}

	body = append(body, jen.Return(
		jen.Op("&").Id(d.TargetMethodExpecterValueStruct).Values(
			jen.Id("expect").Op(":").Id(receiver).Dot("expect"),
		),
	))

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodExpecterStruct)).
		Id("With").Params(params...).Op("*").Id(d.TargetMethodExpecterValueStruct).
		Block(body...).Line()
}

func (d *MethodExpecterData) argumentFuncCode(receiver string, nextStruct string, fn func(string) jen.Code) []jen.Code {
	if len(d.Arguments) == 0 {
		return nil
	}

	var code []jen.Code
	code = append(code, fn(receiver))

	for _, arg := range d.Arguments {
		matchReturnCode := jen.Op("&").Id(nextStruct).Values(
			jen.Id("expect").Op(":").Id("e").Dot("expect"),
			jen.Id("target").Op(":").Id("e").Dot("target"),
		)
		code = append(code, targetMethodExpecterMatchArgCode(
			receiver, nextStruct, matchReturnCode, nextStruct, &d.Lib, arg, false,
		))
	}
	return code
}

func (d *MethodExpecterData) GenerateCode() []jen.Code {
	if d.SkipExpect || (len(d.Arguments) == 0 && len(d.Returns) == 0) {
		return nil
	}

	nm := genlib.NewNameManager("e", nil)
	for _, v := range d.Returns {
		nm.Reserve(v.Name)
	}
	for _, v := range d.Arguments {
		nm.Reserve(v.Name)
	}
	receiver := nm.Request("e")

	code := []jen.Code{
		d.structCode(),
	}

	if len(d.Returns) > 0 {
		code = append(code, targetMethodExpecterReturnCode(
			receiver, d.TargetMethodExpecterStruct, d.TargetMethodReturnStruct, d.Returns,
		))
	}

	code = append(code, d.argumentFuncCode(receiver, d.TargetMethodExpecterMatchArgStruct, d.matchFuncCode)...)
	code = append(code, d.argumentFuncCode(receiver, d.TargetMethodExpecterValueArgStruct, d.withFuncCode)...)
	return code
}
