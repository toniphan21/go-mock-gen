package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

type MethodExpecterValueData struct {
	TargetMethodExpecterValueStruct string
	TargetMethodExpectStruct        string
	TargetMethodReturnStruct        string
	Returns                         []VarInfo
	SkipExpect                      bool
}

func (d *MethodExpecterValueData) structCode() jen.Code {
	return jen.Type().Id(d.TargetMethodExpecterValueStruct).Struct(
		jen.Id("expect").Op("*").Id(d.TargetMethodExpectStruct),
	)
}

func (d *MethodExpecterValueData) GenerateCode() []jen.Code {
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
		targetMethodExpecterReturnCode(receiver, d.TargetMethodExpecterValueStruct, d.TargetMethodReturnStruct, d.Returns),
	}
}
