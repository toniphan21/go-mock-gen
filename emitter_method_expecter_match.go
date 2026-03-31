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

type MethodExpecterMatchData struct {
	TargetMethodExpecterMatchStruct string
	TargetMethodExpectStruct        string
	TargetMethodReturnStruct        string
	Returns                         []VarInfo
	SkipExpect                      bool
}

func (d *MethodExpecterMatchData) structCode() jen.Code {
	return jen.Type().Id(d.TargetMethodExpecterMatchStruct).Struct(
		jen.Id("expect").Op("*").Id(d.TargetMethodExpectStruct),
	)
}

func (d *MethodExpecterMatchData) GenerateCode() []jen.Code {
	if d.SkipExpect || len(d.Returns) == 0 {
		return nil
	}

	nm := genlib.NewNameManager("e", nil)
	for _, v := range d.Returns {
		nm.Reserve(v.Name)
	}
	receiver := nm.Request("e")

	return []jen.Code{
		d.structCode(),
		targetMethodExpecterReturnCode(receiver, d.TargetMethodExpecterMatchStruct, d.TargetMethodReturnStruct, d.Returns),
	}
}
