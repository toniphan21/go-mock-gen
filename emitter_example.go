package mockgen

import (
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

type ExampleData struct {
	Constructor   string
	InterfaceName string
	MethodName    string
	Arguments     []VarInfo
	Returns       []VarInfo
	SkipExpect    bool
	varMock       string
	varSpy        string
}

func (d *ExampleData) argumentsAndReturnsCode() (args, argIds, returns, returnIds, rets []jen.Code) {
	for _, v := range d.Arguments {
		args = append(args, jen.Var().Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		argIds = append(argIds, jen.Id(v.Name))
	}

	for i, v := range d.Returns {
		returns = append(returns, jen.Var().Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		returnIds = append(returnIds, jen.Id(v.Name))
		rets = append(rets, jen.Id("ret"+strconv.Itoa(i)))
	}

	return
}

func (d *ExampleData) testCaseCode(name string, body []jen.Code) jen.Code {
	return jen.Id("t").Dot("Run").Call(
		jen.Lit(name),
		jen.Func().Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(body...),
	)
}

func (d *ExampleData) expectCalledCode(name string, n int) jen.Code {
	if d.SkipExpect {
		return nil
	}

	var args, argIds []jen.Code
	for _, v := range d.Arguments {
		args = append(args, jen.Var().Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		argIds = append(argIds, jen.Id(v.Name))
	}

	var body []jen.Code
	for i := 0; i < n; i++ {
		body = append(body, jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")))
	}
	body = append(body, jen.Line())
	body = append(body, args...)
	for i := 0; i < n; i++ {
		body = append(body, jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...))
	}
	return d.testCaseCode(name, body)
}

func (d *ExampleData) expectCalledStubReturnCode() jen.Code {
	if d.SkipExpect || len(d.Returns) == 0 {
		return nil
	}

	var args, argIds, returns, returnIds, rets = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, returns...)
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("Return").Call(returnIds...),
		jen.Line(),
	)
	body = append(body, args...)
	body = append(body, jen.List(rets...).Op(":=").Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	for i, _ := range d.Returns {
		body = append(body, jen.Qual("fmt", "Println").Call(rets[i], returnIds[i]))
	}
	return d.testCaseCode("expect called - stub return", body)
}

func (d *ExampleData) expectAllUseValue() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
		return nil
	}

	var args, argIds, _, _, _ = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, args...)
	body = append(body, jen.Line())

	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("With").Call(argIds...),
		jen.Line(),
	)
	body = append(body, jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	return d.testCaseCode("expect called - match all arguments by values", body)
}

func (d *ExampleData) expectAllUseValueStubReturn() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 || len(d.Returns) == 0 {
		return nil
	}

	var args, argIds, returns, returnIds, rets = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, args...)
	body = append(body, returns...)
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("With").Call(argIds...).
			Dot("Return").Call(returnIds...),
		jen.Line(),
	)
	body = append(body, jen.List(rets...).Op(":=").Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	for i, _ := range d.Returns {
		body = append(body, jen.Qual("fmt", "Println").Call(rets[i], returnIds[i]))
	}

	return d.testCaseCode("expect called - match all arguments by values - stub return", body)
}

func (d *ExampleData) expectPartialUseValue() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
		return nil
	}

	var args, argIds, _, _, _ = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, args[0])
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("With"+toPascalCase(d.Arguments[0].Name)).Call(argIds[0]),
		jen.Line(),
	)

	for i, v := range args {
		if i != 0 {
			body = append(body, v)
		}
	}
	body = append(body, jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	return d.testCaseCode("expect called - match partial argument by value", body)
}

func (d *ExampleData) expectPartialUseValueStubReturn() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 || len(d.Returns) == 0 {
		return nil
	}

	var args, argIds, returns, returnIds, rets = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, args[0])
	body = append(body, returns...)
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("With"+toPascalCase(d.Arguments[0].Name)).Call(argIds[0]).
			Dot("Return").Call(returnIds...),
		jen.Line(),
	)

	for i, v := range args {
		if i != 0 {
			body = append(body, v)
		}
	}
	body = append(body, jen.List(rets...).Op(":=").Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	for i, _ := range d.Returns {
		body = append(body, jen.Qual("fmt", "Println").Call(rets[i], returnIds[i]))
	}

	return d.testCaseCode("expect called - match partial by value - stub return", body)
}

func (d *ExampleData) expectAllUseCallback() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
		return nil
	}

	var args, argIds, _, _, _ = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("Match").Call(jen.Add(targetMethodMatcherSignature(d.Arguments...)).Add(jen.Block(jen.Return(jen.Lit(true))))),
		jen.Line(),
	)
	body = append(body, args...)
	body = append(body, jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	return d.testCaseCode("expect called - match all arguments by callback", body)
}

func (d *ExampleData) expectAllUseCallbackStubReturn() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 || len(d.Returns) == 0 {
		return nil
	}

	var args, argIds, returns, returnIds, rets = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, returns...)
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("Match").Call(jen.Add(targetMethodMatcherSignature(d.Arguments...)).Add(jen.Block(jen.Return(jen.Lit(true))))).
			Dot("Return").Call(returnIds...),
		jen.Line(),
	)

	body = append(body, args...)
	body = append(body, jen.List(rets...).Op(":=").Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	for i, _ := range d.Returns {
		body = append(body, jen.Qual("fmt", "Println").Call(rets[i], returnIds[i]))
	}

	return d.testCaseCode("expect called - match all arguments by callback - stub return", body)
}

func (d *ExampleData) expectPartialUseCallback() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 {
		return nil
	}

	var args, argIds, _, _, _ = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("Match"+toPascalCase(d.Arguments[0].Name)).
			Call(jen.Add(targetMethodMatcherSignature(d.Arguments[0])).Add(jen.Block(jen.Return(jen.Lit(true))))),
		jen.Line(),
	)

	body = append(body, args...)
	body = append(body, jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	return d.testCaseCode("expect called - match partial argument by callback", body)
}

func (d *ExampleData) expectPartialUseCallbackStubReturn() jen.Code {
	if d.SkipExpect || len(d.Arguments) == 0 || len(d.Returns) == 0 {
		return nil
	}

	var args, argIds, returns, returnIds, rets = d.argumentsAndReturnsCode()

	var body []jen.Code
	body = append(body, returns...)
	body = append(body, jen.Line())
	body = append(body,
		jen.Id(d.varMock).Dot("EXPECT").Call().Dot(d.MethodName).Call(jen.Id("t")).
			Dot("Match"+toPascalCase(d.Arguments[0].Name)).
			Call(jen.Add(targetMethodMatcherSignature(d.Arguments[0])).Add(jen.Block(jen.Return(jen.Lit(true))))).
			Dot("Return").Call(returnIds...),
		jen.Line(),
	)

	body = append(body, args...)
	body = append(body, jen.List(rets...).Op(":=").Id(d.varMock).Dot(d.MethodName).Call(argIds...))

	for i, _ := range d.Returns {
		body = append(body, jen.Qual("fmt", "Println").Call(rets[i], returnIds[i]))
	}
	return d.testCaseCode("expect called - match partial by callback - stub return", body)
}

func (d *ExampleData) stubCode() jen.Code {
	var args, argIds, params, result []jen.Code
	for _, v := range d.Arguments {
		args = append(args, jen.Var().Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
		argIds = append(argIds, jen.Id(v.Name))
		params = append(params, jen.Id(v.Name).Add(genlib.TypeToJenCode(v.Type)))
	}
	for _, v := range d.Returns {
		if v.OriginalName != "" {
			result = append(result, jen.Id(v.OriginalName).Add(genlib.TypeToJenCode(v.Type)))
		} else {
			result = append(result, genlib.TypeToJenCode(v.Type))
		}
	}

	var returned jen.Code
	if len(result) == 0 {
		returned = jen.Line()
	} else {
		var zeros []jen.Code
		for _, v := range d.Returns {
			zeros = append(zeros, genlib.ZeroValueOfType(v.Type))
		}
		returned = jen.Return(zeros...)
	}

	body := []jen.Code{
		jen.Id(d.varMock).Op("=").Id(d.Constructor).Call(),
	}

	body = append(body,
		jen.Id(d.varSpy).Op(":=").Id(d.varMock).Dot("STUB").Call().Dot(d.MethodName).Call(
			jen.Func().Params(params...).Params(result...).Block(returned),
		),
		jen.Line(),
	)
	body = append(body, args...)
	body = append(body,
		jen.Id(d.varMock).Dot(d.MethodName).Call(argIds...),
		jen.Line(),
		jen.Qual("fmt", "Println").Call(jen.Id(d.varSpy)),
	)

	return d.testCaseCode("fine-grained control with stub signature", body)
}

func (d *ExampleData) GenerateCode() []jen.Code {
	nm := genlib.NewNameManager("v", nil)
	for _, v := range d.Returns {
		if v.OriginalName != "" {
			nm.Request(v.OriginalName)
		}
	}
	for _, v := range d.Arguments {
		nm.Request(v.Name)
	}
	d.varMock = nm.Request("mock")
	d.varSpy = nm.Request("spy")

	body := []jen.Code{
		jen.Id(d.varMock).Op(":=").Id(d.Constructor).Call(),
		jen.Line(),
	}

	tests := []jen.Code{
		d.expectCalledCode("expect called once", 1),
		d.expectCalledCode("expect called twice", 2),
		d.expectCalledStubReturnCode(),
		d.expectAllUseValue(),
		d.expectAllUseValueStubReturn(),
		d.expectPartialUseValue(),
		d.expectPartialUseValueStubReturn(),
		d.expectAllUseCallback(),
		d.expectAllUseCallbackStubReturn(),
		d.expectPartialUseCallback(),
		d.expectPartialUseCallbackStubReturn(),
		d.stubCode(),
	}
	for i, v := range tests {
		if v != nil {
			body = append(body, v)
			if i != len(tests)-1 {
				body = append(body, jen.Line())
			}
		}
	}

	return []jen.Code{
		jen.Func().Id(fmt.Sprintf("Test_%s_%s", d.InterfaceName, d.MethodName)).
			Params(jen.Id("t").Op("*").Qual("testing", "T")).
			Block(body...).Line(),
	}
}
