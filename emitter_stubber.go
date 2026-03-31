package mockgen

import (
	"github.com/dave/jennifer/jen"
)

type TargetStubberData struct {
	TargetStruct           string
	TargetStubberStruct    string
	TargetTestDoubleStruct string
	Methods                []MethodInfo
	Lib                    LibraryData
	SkipExpect             bool
}

func (d *TargetStubberData) targetStubberStructCode() jen.Code {
	return jen.Type().Id(d.TargetStubberStruct).Struct(
		jen.Id("target").Op("*").Id(d.TargetStruct),
	).Line()
}

func (d *TargetStubberData) stubCode(receiver string, method MethodInfo) jen.Code {
	vSpy := "m"

	body := []jen.Code{
		jen.If(jen.Id(receiver).Dot("target").Dot("td").Op("==").Nil()).
			Block(jen.Id(receiver).Dot("target").Dot("td").Op("=").
				Op("&").Id(d.TargetTestDoubleStruct).Values()).Line(),

		jen.Id(vSpy).Op(":=").Id(receiver).Dot("target").Dot("td").Dot(method.Name),

		jen.If(jen.Id(vSpy).Op("==").Nil()).Block(
			jen.Id(vSpy).Op("=").Op("&").Id(method.Struct).Values(
				jen.Id("stubLocation").Op(":").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2)),
			),
			jen.Id(receiver).Dot("target").Dot("td").Dot(method.Name).Op("=").Id(vSpy),
		).Line(),

		jen.If(jen.Id("stub").Op("==").Nil()).Block(
			jen.Id(vSpy).Dot("panic").Call(
				jen.Id(d.Lib.MessageStubByNilFunc).Call(
					jen.Id(vSpy),
					jen.Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2)),
				),
			),
		).Line(),

		jen.If(jen.Id(vSpy).Dot("stub").Op("!=").Nil()).Block(
			jen.Id(vSpy).Dot("panic").Call(
				jen.Id(d.Lib.MessageDuplicateStubFunc).Call(
					jen.Id(vSpy),
					jen.Id(vSpy).Dot("stubLocation"),
				),
			),
		).Line(),
	}

	if !d.SkipExpect {
		body = append(body,
			jen.If(jen.Len(jen.Id(vSpy).Dot("expects")).Op(">").Lit(0)).Block(
				jen.Id(vSpy).Dot("panic").Call(
					jen.Id(d.Lib.MessageStubAfterExpectFunc).Call(
						jen.Id(vSpy),
						jen.Id(vSpy).Dot("expects").Index(jen.Lit(0)).Dot("location"),
					),
				),
			).Line(),
		)
	}

	body = append(body,
		jen.Id(vSpy).Dot("stub").Op("=").Id("stub"),
		jen.Return(jen.Id(vSpy)),
	)

	return jen.Func().
		Params(jen.Id(receiver).Op("*").Id(d.TargetStubberStruct)).
		Id(method.Name).
		Params(jen.Id("stub").Add(targetMethodSignature(method.Arguments, method.Returns))).
		Op("*").Id(method.Struct).
		Block(body...)
}

func (d *TargetStubberData) GenerateCode() []jen.Code {
	if len(d.Methods) == 0 {
		return nil
	}

	code := []jen.Code{
		d.targetStubberStructCode(),
	}
	for _, method := range d.Methods {
		code = append(code, d.stubCode("s", method))
	}
	return code
}
