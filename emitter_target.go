package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

type TargetData struct {
	Interface        string
	Struct           string
	Constructor      string
	TestDoubleStruct string
	StubberStruct    string
	ExpecterStruct   string
	Methods          []MethodInfo
	Lib              LibraryData
	SkipExpect       bool
}

func (d *TargetData) constructorCode(location string) jen.Code {
	if d.Constructor == "" {
		return nil
	}
	return jen.Func().Id(d.Constructor).Params().Op("*").Id(d.Struct).Block(
		jen.Return(
			jen.Op("&").Id(d.Struct).Values(
				jen.Id("td").Op(":").Op("&").Id(d.TestDoubleStruct).Values(
					jen.Id(location).Op(":").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(2)),
				),
			),
		),
	)
}

func (d *TargetData) targetStructCode() jen.Code {
	return jen.Type().Id(d.Struct).Struct(
		jen.Id("td").Op("*").Id(d.TestDoubleStruct),
	).Line()
}

func (d *TargetData) testDoubleStructCode(location string) jen.Code {
	var fields = []jen.Code{
		jen.Id(location).String(),
	}

	for _, v := range d.Methods {
		fields = append(fields, jen.Id(v.Name).Op("*").Id(v.Struct))
	}

	return jen.Type().Id(d.TestDoubleStruct).Struct(fields...).Line()
}

func (d *TargetData) targetBuiltinFuncCode(receiver, method, returnedType string) jen.Code {
	return jen.Func().Params(
		jen.Id(receiver).Op("*").Id(d.Struct),
	).Id(method).Params().Op("*").Id(returnedType).Block(
		jen.Return(
			jen.Op("&").Id(returnedType).Values(
				jen.Id("target").Op(":").Id("m"),
			),
		),
	).Line()
}

func (d *TargetData) implementationCode(receiver, location string, method MethodInfo) jen.Code {
	var params, results []jen.Code

	nm := genlib.NewNameManager("v", nil)
	nm.Reserve(receiver)
	for _, v := range method.Arguments {
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		nm.Reserve(v.Name)
	}
	for _, v := range method.Returns {
		if v.OriginalName != "" {
			results = append(results, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
			nm.Reserve(v.OriginalName)
		} else {
			results = append(results, genlib.TypeToJenCode(v.Type))
		}
	}
	vInterface := nm.Next()
	vMethodName := nm.Next()
	vSignature := nm.Next()
	vArgs := nm.Next()
	vMock := nm.Next()
	vIndex := nm.Next()

	var args []jen.Code
	var argIds []jen.Code
	for _, arg := range method.Arguments {
		args = append(args, jen.Lit(arg.Name))
		args = append(args, jen.Id(arg.Name))
		argIds = append(argIds, jen.Id(arg.Name))
	}

	var body []jen.Code

	body = append(body, jen.List(jen.Id(vInterface), jen.Id(vMethodName), jen.Id(vSignature)).Op(":=").List(
		jen.Lit(d.Interface), jen.Lit(method.Name), jen.Lit(targetMethodSignatureString(method)),
	))

	body = append(body, jen.Id(vArgs).Op(":=").Index().Any().Values(args...).Line())

	body = append(body, jen.If(jen.Id(receiver).Dot("td").Op("==").Nil()).Block(
		jen.Panic(jen.Id(d.Lib.MessageNotImplementedFunc).Call(
			jen.Id(vInterface), jen.Id(vMethodName), jen.Id(vSignature), jen.Lit(""), jen.Id(vArgs),
		)),
	).Line())

	switchBody := []jen.Code{}

	var stubCaseBody jen.Code
	if len(method.Returns) > 0 {
		stubCaseBody = jen.Return(jen.Id(vMock).Dot("invokeStub").Call(argIds...))
	} else {
		stubCaseBody = jen.Id(vMock).Dot("invokeStub").Call(argIds...).Line().Return()
	}

	switchBody = append(switchBody, jen.Case(jen.Id(vMock).Dot("stub").Op("!=").Nil()).Block(stubCaseBody))

	if !d.SkipExpect {
		switchBody = append(switchBody, jen.Case(jen.Len(jen.Id(vMock).Dot("expects")).Op(">").Lit(0)).Block(
			jen.Id(vIndex).Op(":=").Len(jen.Id(vMock).Dot("Calls")),
			jen.If(
				jen.Id(vIndex).Op("<").Len(jen.Id(vMock).Dot("expects")),
			).Block(
				jen.Id(vMock).Dot("expects").Index(jen.Id(vIndex)).Dot("tb").Dot("Helper").Call(),
			),
		))

		if len(method.Returns) > 0 {
			switchBody = append(switchBody, jen.Return(jen.Id(vMock).Dot("invokeExpect").Call(argIds...)))
		} else {
			switchBody = append(switchBody, jen.Id(vMock).Dot("invokeExpect").Call(argIds...).Line().Return())
		}
	}

	body = append(body, jen.If(
		jen.Id(vMock).Op(":=").Id(receiver).Dot("td").Dot(method.Name),
		jen.Id(vMock).Op("!=").Nil(),
	).Block(jen.Switch().Block(switchBody...)))

	body = append(body, jen.Panic(jen.Id(d.Lib.MessageNotImplementedFunc).Call(
		jen.Id(vInterface), jen.Id(vMethodName), jen.Id(vSignature),
		jen.Id(receiver).Dot("td").Dot(location),
		jen.Id(vArgs),
	)))

	return jen.Func().
		Params(jen.Id(receiver).Op("*").Id(d.Struct)).
		Id(method.Name).
		Params(params...).Params(results...).Block(body...).
		Line()
}

func (d *TargetData) GenerateCode() []jen.Code {
	if len(d.Methods) == 0 {
		return nil
	}
	nm := genlib.NewNameManager("m", nil)
	for _, method := range d.Methods {
		for _, v := range method.Returns {
			nm.Reserve(v.Name)
		}
		for _, v := range method.Arguments {
			nm.Reserve(v.Name)
		}
	}
	receiver := nm.Request("m")

	lnm := genlib.NewNameManager("l", nil)
	for _, method := range d.Methods {
		lnm.Reserve(method.Name)
	}
	location := lnm.Request("location")

	var code []jen.Code
	if v := d.constructorCode(location); v != nil {
		code = append(code, v)
	}

	code = append(code,
		d.testDoubleStructCode(location),
		d.targetStructCode(),
		d.targetBuiltinFuncCode(receiver, "STUB", d.StubberStruct),
	)

	if !d.SkipExpect {
		code = append(code, d.targetBuiltinFuncCode(receiver, "EXPECT", d.ExpecterStruct))
	}

	for _, method := range d.Methods {
		code = append(code, d.implementationCode(receiver, location, method))
	}

	return code
}
