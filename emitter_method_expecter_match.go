package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

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
