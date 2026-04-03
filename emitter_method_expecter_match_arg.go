package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

func targetMethodExpecterMatchArgCode(receiverName, receiverType string, returnedCode jen.Code, returnedType string, lib *LibraryData, info VarInfo, checkNil bool) jen.Code {
	fn := "Match" + toPascalCase(info.Name)

	param := jen.Id("matcher").Func().Params(jen.Id(info.Name).Add(genlib.TypeToJenCode(info.Type))).Bool()

	body := []jen.Code{
		jen.If(jen.Id("matcher").Op("==").Nil()).Block(
			jen.Id(receiverName).Dot("expect").Dot("tb").Dot("Helper").Call(),
			jen.Id(receiverName).Dot("target").Dot("fatal").Call(
				jen.Id(receiverName).Dot("expect").Dot("index"),
				jen.Id(lib.MessageMatchArgByNilFunc).Call(
					jen.Id(receiverName).Dot("target"),
					jen.Lit(fn),
				),
			),
		),
	}

	if checkNil {
		body = append(body,
			jen.Line(),
			jen.If(jen.Id(receiverName).Dot("expect").Dot("matcher").Dot(info.Name).Op("!=").Nil()).Block(
				jen.Id(receiverName).Dot("expect").Dot("tb").Dot("Helper").Call(),
				jen.Id(receiverName).Dot("target").Dot("fatal").Call(
					jen.Id(receiverName).Dot("expect").Dot("index"),
					jen.Id(lib.MessageDuplicateMatchArgFunc).Call(
						jen.Id(receiverName).Dot("target"),
						jen.Lit(fn),
						jen.Id(receiverName).Dot("expect").Dot("matcherLocations").Index(jen.Lit(info.Name)),
					),
				),
			),
		)
	}

	body = append(body,
		jen.Line(),
		jen.Id(receiverName).Dot("expect").Dot("matcher").Dot(info.Name).Op("=").Id("matcher"),
		jen.Id(receiverName).Dot("expect").Dot("matcherLocations").Index(jen.Lit(info.Name)).Op("=").Id(lib.CallerLocationFunc).Call(jen.Lit(2)),
		jen.Id(receiverName).Dot("expect").Dot("matcherHints").Index(jen.Lit(info.Name)).Op("=").Id(lib.MessageMatchArgHintFunc).Call(),
		jen.Return(returnedCode),
	)

	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id(fn).Params(param).Op("*").Id(returnedType).
		Block(body...).
		Line()
}

type MethodExpecterMatchArgData struct {
	ExpecterMatchArgStruct string
	ExpectStruct           string
	Struct                 string
	ReturnStruct           string
	Arguments              []VarInfo
	Returns                []VarInfo
	Lib                    LibraryData
	SkipExpect             bool
}

func (d *MethodExpecterMatchArgData) structCode() jen.Code {
	return jen.Type().Id(d.ExpecterMatchArgStruct).Struct(
		jen.Id("expect").Op("*").Id(d.ExpectStruct),
		jen.Id("target").Op("*").Id(d.Struct),
	).Line()
}

func (d *MethodExpecterMatchArgData) GenerateCode() []jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
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
			receiver, d.ExpecterMatchArgStruct, d.ReturnStruct, d.Returns,
		))
	}

	for _, arg := range d.Arguments {
		code = append(code, targetMethodExpecterMatchArgCode(
			receiver, d.ExpecterMatchArgStruct, jen.Id(receiver), d.ExpecterMatchArgStruct, &d.Lib, arg, true,
		))
	}
	return code
}
