package mockgen

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

func targetMethodExpecterValueArgCode(receiverName, receiverType string, returnedCode jen.Code, returnedType string, lib *LibraryData, info VarInfo, checkNil bool) jen.Code {
	fn := "With" + toPascalCase(info.Name)
	matchFunc := lib.ReflectEqualMatcherFunc
	matchMethod := "reflect.DeepEqual"

	if types.Comparable(info.Type) {
		matchFunc = lib.BasicComparisonMatcherFunc
		matchMethod = "=="
	}

	var body []jen.Code
	if checkNil {
		body = append(body,
			jen.If(
				jen.Id(receiverName).Dot("expect").Dot("matcher").Dot(info.Name).Op("!=").Nil(),
			).Block(
				jen.Id(receiverName).Dot("expect").Dot("tb").Dot("Helper").Call(),
				jen.Id(receiverName).Dot("target").Dot("fatal").Call(
					jen.Id(receiverName).Dot("expect").Dot("index"),
					jen.Id(lib.MessageDuplicateMatchArgFunc).Call(
						jen.Id(receiverName).Dot("target"),
						jen.Lit(fn),
						jen.Id(receiverName).Dot("expect").Dot("matcherLocations").Index(jen.Lit(info.Name)),
					),
				),
			).Line(),
		)
	}

	body = append(body,
		jen.Id(receiverName).Dot("expect").Dot("matcher").Dot(info.Name).Op("=").Id(matchFunc).Call(jen.Id(info.Name)),
		jen.Id(receiverName).Dot("expect").Dot("matcherWants").Index(jen.Lit(info.Name)).Op("=").Id(info.Name),
		jen.Id(receiverName).Dot("expect").Dot("matcherMethods").Index(jen.Lit(info.Name)).Op("=").Lit(matchMethod),
		jen.Id(receiverName).Dot("expect").Dot("matcherLocations").Index(jen.Lit(info.Name)).Op("=").Id(lib.CallerLocationFunc).Call(jen.Lit(2)),
		jen.Line(),
		jen.Return(returnedCode),
	)

	return jen.Func().Params(jen.Id(receiverName).Op("*").Id(receiverType)).
		Id(fn).Params(jen.Id(info.Name).Add(genlib.TypeToJenCode(info.Type))).Op("*").Id(returnedType).
		Block(body...).
		Line()
}

type MethodExpecterValueArgData struct {
	ExpecterValueArgStruct string
	ExpectStruct           string
	Struct                 string
	ReturnStruct           string
	Arguments              []VarInfo
	Returns                []VarInfo
	Lib                    LibraryData
	SkipExpect             bool
}

func (d *MethodExpecterValueArgData) structCode() jen.Code {
	return jen.Type().Id(d.ExpecterValueArgStruct).Struct(
		jen.Id("expect").Op("*").Id(d.ExpectStruct),
		jen.Id("target").Op("*").Id(d.Struct),
	).Line()
}

func (d *MethodExpecterValueArgData) GenerateCode() []jen.Code {
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
			receiver, d.ExpecterValueArgStruct, d.ReturnStruct, d.Returns,
		))
	}

	for _, arg := range d.Arguments {
		code = append(code, targetMethodExpecterValueArgCode(
			receiver, d.ExpecterValueArgStruct, jen.Id(receiver), d.ExpecterValueArgStruct, &d.Lib, arg, true,
		))
	}
	return code
}
