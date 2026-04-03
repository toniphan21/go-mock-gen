package mockgen

import "github.com/dave/jennifer/jen"

type TargetExpecterData struct {
	Struct           string
	ExpecterStruct   string
	TestDoubleStruct string
	Methods          []MethodInfo
	Lib              LibraryData
	SkipExpect       bool
}

func (d *TargetExpecterData) targetExpecterStructCode() jen.Code {
	return jen.Type().Id(d.ExpecterStruct).Struct(
		jen.Id("target").Op("*").Id(d.Struct),
	).Line()
}

func (d *TargetExpecterData) expectCode(receiver string, method MethodInfo) jen.Code {
	vMock := "m"
	vIndex := "idx"

	body := []jen.Code{
		jen.If(jen.Id(receiver).Dot("target").Dot("td").Op("==").Nil()).Block(
			jen.Id(receiver).Dot("target").Dot("td").Op("=").Op("&").Id(d.TestDoubleStruct).Values(),
		),
		jen.Line(),
		jen.Var().Id(vMock).Op("=").Id(receiver).Dot("target").Dot("td").Dot(method.Name),
		jen.If(jen.Id(vMock).Op("==").Nil()).Block(
			jen.Id(vMock).Op("=").Op("&").Id(method.Struct).Values(),
			jen.Id(receiver).Dot("target").Dot("td").Dot(method.Name).Op("=").Id(vMock),
		),
		jen.Line(),
		jen.If(jen.Id(vMock).Dot("stub").Op("!=").Nil()).Block(
			jen.Id(vMock).Dot("panic").Call(
				jen.Id(d.Lib.MessageExpectAfterStubFunc).Call(
					jen.Id(vMock),
					jen.Id(vMock).Dot("stubLocation"),
				),
			),
		),
		jen.Line(),
		jen.If(jen.Id("tb").Op("==").Nil()).Block(
			jen.Id(vMock).Dot("panic").Call(
				jen.Id(d.Lib.MessageExpectByNilFunc).Call(jen.Id(vMock)),
			),
		),
		jen.Line(),
		jen.Id(vIndex).Op(":=").Len(jen.Id(vMock).Dot("expects")),
		jen.Id(vMock).Dot("expects").Op("=").Append(
			jen.Id(vMock).Dot("expects"),
			jen.Op("&").Id(method.ExpectStruct).ValuesFunc(func(g *jen.Group) {
				g.Line().Id("location").Op(":").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2))
				if len(method.Arguments) > 0 {
					g.Line().Id("matcher").Op(":").Op("&").Id(method.ArgumentMatcherStruct).Values()
					g.Line().Id("matcherWants").Op(":").Make(jen.Map(jen.String()).Any())
					g.Line().Id("matcherMethods").Op(":").Make(jen.Map(jen.String()).String())
					g.Line().Id("matcherHints").Op(":").Make(jen.Map(jen.String()).String())
					g.Line().Id("matcherLocations").Op(":").Make(jen.Map(jen.String()).String())
				}
				g.Line().Id("index").Op(":").Id(vIndex)
				g.Line().Id("tb").Op(":").Id("tb")
				g.Line()
			}),
		),
		jen.Line(),
		jen.Id("tb").Dot("Helper").Call(),
		jen.Id("tb").Dot("Cleanup").Call(
			jen.Func().Params().Block(
				jen.Id("tb").Dot("Helper").Call(),
				jen.Id(vMock).Dot("verify").Call(jen.Id(vIndex)),
			),
		),
	}

	if len(method.Returns) > 0 {
		var returnCode jen.Code
		if len(method.Arguments) > 0 {
			returnCode = jen.Return(jen.Op("&").Id(method.ExpecterStruct).Values(
				jen.Id("target").Op(":").Id(vMock),
				jen.Id("expect").Op(":").Id(vMock).Dot("expects").Index(jen.Id(vIndex)),
			))
		} else {
			returnCode = jen.Return(jen.Op("&").Id(method.ExpecterStruct).Values(
				jen.Id("expect").Op(":").Id(vMock).Dot("expects").Index(jen.Id(vIndex)),
			))
		}

		body = append(body, jen.Line(), returnCode)

		return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.ExpecterStruct)).
			Id(method.Name).Params(jen.Id("tb").Qual("testing", "TB")).
			Op("*").Id(method.ExpecterStruct).
			Block(body...).Line()
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.ExpecterStruct)).
		Id(method.Name).Params(jen.Id("tb").Qual("testing", "TB")).
		Block(body...).Line()
}

func (d *TargetExpecterData) GenerateCode() []jen.Code {
	if d.SkipExpect || len(d.Methods) == 0 {
		return nil
	}

	code := []jen.Code{
		d.targetExpecterStructCode(),
	}
	for _, method := range d.Methods {
		code = append(code, d.expectCode("e", method))
	}
	return code
}
