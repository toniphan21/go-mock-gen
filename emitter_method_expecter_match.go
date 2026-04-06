package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

type MethodExpecterMatchData struct {
	ExpecterMatchStruct string
	ExpectStruct        string
	ReturnStruct        string
	Arguments           []VarInfo
	Returns             []VarInfo
	SkipExpect          bool
}

func (d *MethodExpecterMatchData) structCode() jen.Code {
	return jen.Type().Id(d.ExpecterMatchStruct).Struct(
		jen.Id("expect").Op("*").Id(d.ExpectStruct),
	)
}

func (d *MethodExpecterMatchData) GenerateCode() []jen.Code {
	if d.SkipExpect || len(d.Returns) == 0 || len(d.Arguments) == 0 {
		return nil
	}

	nm := genlib.NewNameManager("e", nil)
	for _, v := range d.Returns {
		nm.Reserve(v.Name)
	}
	receiver := nm.Request("e")

	return []jen.Code{
		d.structCode(),
		targetMethodExpecterReturnCode(receiver, d.ExpecterMatchStruct, d.ReturnStruct, d.Returns),
	}
}
