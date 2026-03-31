package mockgen

import (
	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

//	func (m *targetFull) methodName() string {
//		return "Full"
//	}
//
//	func (m *targetFull) interfaceName() string {
//		return "Target"
//	}
//
//	func (m *targetFull) fatal(index int, msg string) {
//		m.verified = true              // skip:!expect
//		m.expects[index].tb.Helper()   // skip:!expect
//		m.expects[index].tb.Fatal(msg) // skip:!expect
//	}
//
//	func (m *targetFull) panic(msg string) {
//		m.verified = true // skip:!expect
//		panic(msg)
//	}
//
//	func (m *targetFull) buildCallHistory(sb *strings.Builder, header string) {
//		if header != "" && len(m.Calls) != 0 { // skip:!expect
//			sb.WriteString(fmt.Sprintf("%s:\n", header))
//		}
//
//		for i, call := range m.Calls { // skip:!expect
//			args := []any{"ctx", call.Argument.ctx, "input", call.Argument.input}
//			libMessageCallHistory(sb, i, m.expects[i].location, call.Location, args)
//		}
//	}
//
//	func (m *targetFull) invokeStub(ctx context.Context, input string) ([]Result, error) {
//		v0, v1 := m.stub(ctx, input)
//		return m.capture(
//			targetFullArgument{ctx: ctx, input: input},
//			targetFullReturn{first: v0, second: v1},
//		)
//	}
//
// func (m *targetFull) invokeExpect(ctx context.Context, input string) ([]Result, error) { // skip:!expect
//
//		args := []any{"ctx", ctx, "input", input}
//		index := len(m.Calls)
//		if index >= len(m.expects) {
//			m.panic(libMessageTooManyCalls(m, len(m.expects), index+1, args))
//		}
//
//		expect := m.expects[index]
//		if expect.match != nil && !expect.match(ctx, input) {
//			expect.tb.Helper()
//			m.fatal(index, libMessageMatchFail(m, expect.matchLocation, index, args))
//		}
//
//		expect.tb.Helper()
//		libMatchArgument(m, index, "ctx", ctx, expect.matcher.ctx, expect.matcherWants, expect.matcherMethods, expect.matcherHints, expect.tb, expect.matcherLocations["ctx"])
//		libMatchArgument(m, index, "input", input, expect.matcher.input, expect.matcherWants, expect.matcherMethods, expect.matcherHints, expect.tb, expect.matcherLocations["input"])
//
//		return m.capture(
//			targetFullArgument{ctx: ctx, input: input},
//			expect.returns,
//		)
//	}
//
//	func (m *targetFull) capture(args targetFullArgument, returns targetFullReturn) ([]Result, error) {
//		m.Calls = append(m.Calls, targetFullCall{
//			Location:  libCallerLocation(4),
//			Argument: args,
//			Return:   returns,
//		})
//		return returns.first, returns.second
//	}
//
// func (m *targetFull) verify(index int) { // skip:!expect
//		if !m.verified && index >= len(m.Calls) {
//			m.expects[index].tb.Helper()
//			m.expects[index].tb.Fatal(libMessageExpectButNotCalled(m, len(m.expects), len(m.Calls), index))
//		}
//	}

type MethodData struct {
	TargetMethodStruct                string
	TargetMethodCallStruct            string
	TargetMethodArgumentStruct        string
	TargetMethodArgumentMatcherStruct string
	TargetMethodReturnStruct          string
	TargetMethodExpectStruct          string
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
	code := []jen.Code{
		d.structCode(),
	}
	if v := d.expectStructCode(); v != nil {
		code = append(code, v)
	}
	return code
}
