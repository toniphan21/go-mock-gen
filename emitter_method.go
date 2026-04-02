package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

type MethodData struct {
	TargetMethodStruct                string
	TargetMethodCallStruct            string
	TargetMethodArgumentStruct        string
	TargetMethodArgumentMatcherStruct string
	TargetMethodReturnStruct          string
	TargetMethodExpectStruct          string
	Interface                         string
	Name                              string
	Arguments                         []VarInfo
	Returns                           []VarInfo
	Lib                               LibraryData
	SkipExpect                        bool
}

func (d *MethodData) structCode() jen.Code {
	return jen.Type().Id(d.TargetMethodStruct).StructFunc(func(g *jen.Group) {
		g.Id("Calls").Index().Id(d.TargetMethodCallStruct)
		g.Id("stub").Add(targetMethodSignature(d.Arguments, d.Returns))
		g.Id("stubLocation").String()
		if !d.SkipExpect {
			g.Id("expects").Index().Op("*").Id(d.TargetMethodExpectStruct)
			g.Id("verified").Bool()
		}
	}).Line()
}

func (d *MethodData) methodNameFuncCode(receiver string) jen.Code {
	return jen.Func().
		Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("methodName").Params().String().
		Block(jen.Return(jen.Lit(d.Name))).Line()
}

func (d *MethodData) interfaceNameFuncCode(receiver string) jen.Code {
	return jen.Func().
		Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("interfaceName").Params().String().
		Block(jen.Return(jen.Lit(d.Interface))).Line()
}

func (d *MethodData) fatalFuncCode(receiver string) jen.Code {
	var body []jen.Code
	if !d.SkipExpect {
		body = append(body,
			jen.Id(receiver).Dot("verified").Op("=").Lit(true),
			jen.Id(receiver).Dot("expects").Index(jen.Id("index")).Dot("tb").Dot("Helper").Call(),
			jen.Id(receiver).Dot("expects").Index(jen.Id("index")).Dot("tb").Dot("Fatal").Call(jen.Id("msg")),
		)
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).Id("fatal").
		Params(jen.Id("index").Int(), jen.Id("msg").String()).
		Block(body...).Line()
}

func (d *MethodData) panicFuncCode(receiver string) jen.Code {
	var body []jen.Code
	if !d.SkipExpect {
		body = append(body, jen.Id(receiver).Dot("verified").Op("=").Lit(true))
	}
	body = append(body, jen.Panic(jen.Id("msg")))

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).Id("panic").
		Params(jen.Id("msg").String()).
		Block(body...).Line()
}

func (d *MethodData) buildCallHistoryFuncCode(receiver string) jen.Code {
	var body []jen.Code
	if !d.SkipExpect {
		var args []jen.Code
		for _, arg := range d.Arguments {
			args = append(args, jen.Lit(arg.Name), jen.Id("v").Dot("Argument").Dot(arg.Field))
		}

		argsCode := jen.Id("a").Op(":=").Index().Any().Values(args...)

		body = []jen.Code{
			jen.If(
				jen.Id("header").Op("!=").Lit("").Op("&&").
					Len(jen.Id(receiver).Dot("Calls")).Op("!=").Lit(0),
			).Block(
				jen.Id("sb").Dot("WriteString").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s:\n"), jen.Id("header")),
				),
			),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().
					Id(receiver).Dot("Calls"),
			).Block(
				argsCode,
				jen.Id(d.Lib.MessageCallHistoryFunc).Call(
					jen.Id("sb"),
					jen.Id("i"),
					jen.Id(receiver).Dot("expects").Index(jen.Id("i")).Dot("location"),
					jen.Id("v").Dot("Location"),
					jen.Id("a"),
				),
			),
		}
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).Id("buildCallHistory").
		Params(
			jen.Id("sb").Op("*").Qual("strings", "Builder"),
			jen.Id("header").String(),
		).
		Block(body...).Line()
}

func (d *MethodData) invokeStubFuncCode(receiver string) jen.Code {
	var params, results, argumentFields, returnFields, captureArgs, passed, vars []jen.Code
	var body []jen.Code

	nm := genlib.NewNameManager("v", nil)
	nm.Reserve(receiver)
	for _, v := range d.Arguments {
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		nm.Reserve(v.Name)

		passed = append(passed, jen.Id(v.Name))
		argumentFields = append(argumentFields, jen.Id(v.Field).Op(":").Id(v.Name))
	}

	for _, v := range d.Returns {
		if v.OriginalName != "" {
			results = append(results, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
			nm.Reserve(v.OriginalName)
		} else {
			results = append(results, genlib.TypeToJenCode(v.Type))
		}
	}

	for _, v := range d.Returns {
		vn := nm.Next()
		vars = append(vars, jen.Id(vn))
		returnFields = append(returnFields, jen.Id(v.Field).Op(":").Id(vn))
	}

	if len(d.Arguments) > 0 {
		captureArgs = append(captureArgs, jen.Id(d.TargetMethodArgumentStruct).Values(argumentFields...))
	}

	if len(d.Returns) == 0 {
		body = append(body, jen.Id(receiver).Dot("stub").Call(passed...))
		body = append(body, jen.Id(receiver).Dot("capture").Call(captureArgs...))
	} else {
		captureArgs = append(captureArgs, jen.Id(d.TargetMethodReturnStruct).Values(returnFields...))
		body = append(body, jen.List(vars...).Op(":=").Id(receiver).Dot("stub").Call(passed...))
		body = append(body, jen.Return(jen.Id(receiver).Dot("capture").Call(captureArgs...)))
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("invokeStub").
		Params(params...).Params(results...).
		Block(body...).Line()
}

func (d *MethodData) invokeExpectFuncCode(receiver string) jen.Code {
	if d.SkipExpect {
		return nil
	}

	var params, results, args, argIds []jen.Code

	nm := genlib.NewNameManager("v", nil)
	nm.Reserve(receiver)
	for _, v := range d.Arguments {
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		args = append(args, jen.Lit(v.Name))
		args = append(args, jen.Id(v.Name))
		argIds = append(argIds, jen.Id(v.Name))
		nm.Reserve(v.Name)
	}
	for _, v := range d.Returns {
		if v.OriginalName != "" {
			results = append(results, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
			nm.Reserve(v.OriginalName)
		} else {
			results = append(results, genlib.TypeToJenCode(v.Type))
		}
	}

	vArgs := nm.Next()
	vIndex := nm.Next()
	vExpect := nm.Next()

	body := []jen.Code{
		jen.Id(vArgs).Op(":=").Index().Any().Values(args...).Line(),
		jen.Id(vIndex).Op(":=").Len(jen.Id(receiver).Dot("Calls")),
		jen.If(jen.Id(vIndex).Op(">=").Len(jen.Id(receiver).Dot("expects"))).Block(
			jen.Id(receiver).Dot("panic").Call(
				jen.Id(d.Lib.MessageTooManyCallsFunc).Call(
					jen.Id(receiver),
					jen.Len(jen.Id(receiver).Dot("expects")),
					jen.Id(vIndex).Op("+").Lit(1),
					jen.Id(vArgs),
				),
			),
		).Line(),

		jen.Id(vExpect).Op(":=").Id(receiver).Dot("expects").Index(jen.Id(vIndex)),
		jen.If(
			jen.Id(vExpect).Dot("match").Op("!=").Nil().Op("&&").
				Op("!").Id(vExpect).Dot("match").Call(argIds...),
		).Block(
			jen.Id(vExpect).Dot("tb").Dot("Helper").Call(),
			jen.Id(receiver).Dot("fatal").Call(
				jen.Id(vIndex),
				jen.Id(d.Lib.MessageMatchFailFunc).Call(
					jen.Id(receiver), jen.Id(vExpect).Dot("matchLocation"), jen.Id(vIndex), jen.Id(vArgs),
				),
			),
		).Line(),
	}

	if len(d.Arguments) > 0 {
		body = append(body, jen.Id(vExpect).Dot("tb").Dot("Helper").Call())
	}

	var argFields []jen.Code
	for _, arg := range d.Arguments {
		body = append(body,
			jen.Id(d.Lib.MatchArgumentFunc).Call(
				jen.Id(receiver),
				jen.Id(vIndex),
				jen.Lit(arg.Name),
				jen.Id(arg.Name),
				jen.Id(vExpect).Dot("matcher").Dot(arg.Name),
				jen.Id(vExpect).Dot("matcherWants"),
				jen.Id(vExpect).Dot("matcherMethods"),
				jen.Id(vExpect).Dot("matcherHints"),
				jen.Id(vExpect).Dot("tb"),
				jen.Id(vExpect).Dot("matcherLocations").Index(jen.Lit(arg.Name)),
			),
		)
		argFields = append(argFields, jen.Id(arg.Field).Op(":").Id(arg.Name))
	}

	if len(d.Arguments) > 0 {
		body = append(body, jen.Line())
	}

	var captureArgs []jen.Code
	if len(d.Arguments) > 0 {
		captureArgs = append(captureArgs, jen.Id(d.TargetMethodArgumentStruct).Values(argFields...))
	}

	if len(d.Returns) == 0 {
		body = append(body, jen.Id(receiver).Dot("capture").Call(captureArgs...))
	} else {
		captureArgs = append(captureArgs, jen.Id("expect").Dot("returns"))
		body = append(body, jen.Return(jen.Id(receiver).Dot("capture").Call(captureArgs...)))
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("invokeExpect").Params(params...).Params(results...).
		Block(body...).Line()
}

func (d *MethodData) captureFuncCode(receiver string) jen.Code {
	var params, results, returns, callFields []jen.Code

	for _, v := range d.Returns {
		if v.OriginalName != "" {
			results = append(results, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
		} else {
			results = append(results, genlib.TypeToJenCode(v.Type))
		}
		returns = append(returns, jen.Id("returns").Dot(v.Field))
	}

	callFields = append(callFields, jen.Id("Location").Op(":").Id(d.Lib.CallerLocationFunc).Call(jen.Lit(4)))
	if len(d.Arguments) > 0 {
		params = append(params, jen.Id("args").Id(d.TargetMethodArgumentStruct))
		callFields = append(callFields, jen.Id("Argument").Op(":").Id("args"))
	}

	if len(d.Returns) > 0 {
		params = append(params, jen.Id("returns").Id(d.TargetMethodReturnStruct))
		callFields = append(callFields, jen.Id("Return").Op(":").Id("returns"))
	}

	body := []jen.Code{
		jen.Id(receiver).Dot("Calls").Op("=").Append(
			jen.Id(receiver).Dot("Calls"),
			jen.Id(d.TargetMethodCallStruct).Values(callFields...),
		),
	}

	if len(results) > 0 {
		body = append(body, jen.Return(returns...))
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("capture").Params(params...).Params(results...).
		Block(body...).Line()
}

func (d *MethodData) verifyFuncCode(receiver string) jen.Code {
	if d.SkipExpect {
		return nil
	}

	return jen.Func().Params(jen.Id(receiver).Op("*").Id(d.TargetMethodStruct)).
		Id("verify").Params(jen.Id("index").Int()).
		Block(jen.If(
			jen.Op("!").Id(receiver).Dot("verified").
				Op("&&").
				Id("index").Op(">=").Len(jen.Id(receiver).Dot("Calls")),
		).Block(
			jen.Id(receiver).Dot("expects").Index(jen.Id("index")).Dot("tb").Dot("Helper").Call(),
			jen.Id(receiver).Dot("expects").Index(jen.Id("index")).Dot("tb").Dot("Fatal").Call(
				jen.Id(d.Lib.MessageExpectButNotCalledFunc).Call(
					jen.Id(receiver),
					jen.Len(jen.Id(receiver).Dot("expects")),
					jen.Len(jen.Id(receiver).Dot("Calls")),
					jen.Id("index"),
				),
			),
		)).Line()
}

func (d *MethodData) callStructCode() jen.Code {
	return jen.Type().Id(d.TargetMethodCallStruct).StructFunc(func(g *jen.Group) {
		g.Id("Location").String()

		if len(d.Arguments) > 0 {
			g.Id("Argument").Id(d.TargetMethodArgumentStruct)
		}

		if len(d.Returns) > 0 {
			g.Id("Return").Id(d.TargetMethodReturnStruct)
		}
	}).Line()
}

func (d *MethodData) argumentStructCode() jen.Code {
	if len(d.Arguments) == 0 {
		return nil
	}

	return jen.Type().Id(d.TargetMethodArgumentStruct).StructFunc(func(g *jen.Group) {
		for _, v := range d.Arguments {
			g.Id(v.Field).Add(genlib.TypeToJenCode(v.Type))
		}
	}).Line()
}

func (d *MethodData) argumentMatcherStructCode() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
		return nil
	}

	return jen.Type().Id(d.TargetMethodArgumentMatcherStruct).StructFunc(func(g *jen.Group) {
		for _, v := range d.Arguments {
			g.Id(v.Field).Add(targetMethodMatcherSignature(v))
		}
	}).Line()
}

func (d *MethodData) returnStructCode() jen.Code {
	if len(d.Returns) == 0 {
		return nil
	}

	return jen.Type().Id(d.TargetMethodReturnStruct).StructFunc(func(g *jen.Group) {
		for _, v := range d.Returns {
			g.Id(v.Field).Add(genlib.TypeToJenCode(v.Type))
		}
	}).Line()
}

func (d *MethodData) expectStructCode() jen.Code {
	if d.SkipExpect {
		return nil
	}

	return jen.Type().Id(d.TargetMethodExpectStruct).StructFunc(func(g *jen.Group) {
		if len(d.Arguments) > 0 {
			g.Id("match").Add(targetMethodMatcherSignature(d.Arguments...))
			g.Id("matchLocation").String()
			g.Id("matcher").Op("*").Id(d.TargetMethodArgumentMatcherStruct)
			g.Id("matcherWants").Map(jen.String()).Any()
			g.Id("matcherMethods").Map(jen.String()).String()
			g.Id("matcherHints").Map(jen.String()).String()
			g.Id("matcherLocations").Map(jen.String()).String()
		}

		if len(d.Returns) > 0 {
			g.Id("returns").Id(d.TargetMethodReturnStruct)
		}

		g.Id("location").String()
		g.Id("index").Int()
		g.Id("tb").Qual("testing", "TB")
	}).Line()
}

func (d *MethodData) GenerateCode() []jen.Code {
	nm := genlib.NewNameManager("m", nil)
	for _, v := range d.Arguments {
		nm.Reserve(v.Name)
	}
	for _, v := range d.Returns {
		if v.OriginalName != "" {
			nm.Reserve(v.OriginalName)
		}
	}
	receiver := nm.Request("m")

	parts := []jen.Code{
		d.structCode(),
		d.methodNameFuncCode(receiver),
		d.interfaceNameFuncCode(receiver),
		d.fatalFuncCode(receiver),
		d.panicFuncCode(receiver),
		d.buildCallHistoryFuncCode(receiver),
		d.invokeStubFuncCode(receiver),
		d.invokeExpectFuncCode(receiver),
		d.captureFuncCode(receiver),
		d.verifyFuncCode(receiver),
		d.callStructCode(),
		d.argumentStructCode(),
		d.argumentMatcherStructCode(),
		d.returnStructCode(),
		d.expectStructCode(),
	}

	var code []jen.Code
	for _, v := range parts {
		if v != nil {
			code = append(code, v)
		}
	}
	return code
}
