package mockgen

import "github.com/dave/jennifer/jen"

type LibraryData struct {
	CallerLocationFunc            string
	MethodInterface               string
	MessageWriteArgumentsFunc     string
	MessageMatchFailFunc          string
	MessageArgumentMismatchedFunc string
	MessageNotImplementedFunc     string
	MessageCallHistoryFunc        string
	MessageTooManyCallsFunc       string
	MessageMatchByNilFunc         string
	MessageExpectByNilFunc        string
	MessageExpectAfterStubFunc    string
	MessageStubByNilFunc          string
	MessageStubAfterExpectFunc    string
	MessageDuplicateStubFunc      string
	MessageExpectButNotCalledFunc string
	CompareByReflectEqualFunc     string
	CompareByBasicComparisonFunc  string
}

func (lib *LibraryData) CallerLocationCode() jen.Code {
	return jen.Func().Id(lib.CallerLocationFunc).Params(
		jen.Id("skip").Int(),
	).String().Block(
		jen.List(
			jen.Id("_"),
			jen.Id("file"),
			jen.Id("line"),
			jen.Id("_"),
		).Op(":=").Qual("runtime", "Caller").Call(jen.Id("skip")),
		jen.Return(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s:%d"),
				jen.Qual("path/filepath", "Base").Call(jen.Id("file")),
				jen.Id("line"),
			),
		),
	)
}

func (lib *LibraryData) MethodInterfaceCode() jen.Code {
	return jen.Type().Id(lib.MethodInterface).Interface(
		jen.Id("methodName").Params().String(),
		jen.Id("interfaceName").Params().String(),
		jen.Id("buildCallHistory").Params(
			jen.Id("sb").Op("*").Qual("strings", "Builder"),
			jen.Id("header").String(),
		),
		jen.Id("fatal").Params(
			jen.Id("index").Int(),
			jen.Id("msg").String(),
		),
		jen.Id("panic").Params(
			jen.Id("msg").String(),
		),
	)
}

func (lib *LibraryData) MessageWriteArgumentsCode() jen.Code {
	return jen.Func().Id(lib.MessageWriteArgumentsFunc).Params(
		jen.Id("sb").Op("*").Qual("strings", "Builder"),
		jen.Id("template").String(),
		jen.Id("args").Index().Any(),
	).Block(
		jen.Id("maxLen").Op(":=").Lit(0),
		jen.For(
			jen.Id("i").Op(":=").Lit(0),
			jen.Id("i").Op("<").Len(jen.Id("args")),
			jen.Id("i").Op("+=").Lit(2),
		).Block(
			jen.List(jen.Id("str"), jen.Id("ok")).Op(":=").Id("args").Index(jen.Id("i")).Assert(jen.String()),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Id("str").Op("=").Qual("fmt", "Sprintf").Call(
					jen.Lit("%v"),
					jen.Id("args").Index(jen.Id("i")),
				),
			),
			jen.Id("maxLen").Op("=").Id("max").Call(
				jen.Id("maxLen"),
				jen.Len(jen.Id("str")),
			),
		),
		jen.Line(),
		jen.Id("format").Op(":=").Qual("strings", "ReplaceAll").Call(
			jen.Id("template"),
			jen.Lit("[MAX-KEY-LEN]"),
			jen.Qual("strconv", "Itoa").Call(jen.Id("maxLen")),
		),
		jen.For(
			jen.Id("i").Op(":=").Lit(0),
			jen.Id("i").Op("<").Len(jen.Id("args")),
			jen.Id("i").Op("+=").Lit(2),
		).Block(
			jen.List(jen.Id("key"), jen.Id("ok")).Op(":=").Id("args").Index(jen.Id("i")).Assert(jen.String()),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Id("key").Op("=").Qual("fmt", "Sprintf").Call(
					jen.Lit("%v"),
					jen.Id("args").Index(jen.Id("i")),
				),
			),
			jen.Line(),
			jen.Var().Id("val").Any(),
			jen.If(jen.Id("i").Op("+").Lit(1).Op("<").Len(jen.Id("args"))).Block(
				jen.Id("val").Op("=").Id("args").Index(jen.Id("i").Op("+").Lit(1)),
			),
			jen.Id("sb").Dot("WriteString").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Id("format"),
					jen.Id("key"),
					jen.Id("val"),
				),
			),
		),
	)
}

func (lib *LibraryData) MessageMatchFailCode() jen.Code {
	return jen.Func().Id(lib.MessageMatchFailFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("matchedAt").String(),
		jen.Id("index").Int(),
		jen.Id("args").Index().Any(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s call #%d did not match\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
				jen.Id("index").Op("+").Lit(1),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("arguments:\n"),
			),
		),
		jen.Id(lib.MessageWriteArgumentsFunc).Call(
			jen.Id("sb"),
			jen.Lit("\t%[MAX-KEY-LEN]s = %#v\n"),
			jen.Id("args"),
		),
		jen.Id("sb").Dot("WriteString").Call(jen.Lit("\n")),
		jen.Id("m").Dot("buildCallHistory").Call(
			jen.Id("sb"),
			jen.Lit("call history"),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("hint: check the callback passed to Match at %s"),
				jen.Id("matchedAt"),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageArgumentMismatchedCode() jen.Code {
	return jen.Func().Id(lib.MessageArgumentMismatchedFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("argName").String(),
		jen.Id("expectAt").String(),
		jen.Id("comparedBy").String(),
		jen.Id("callNo").Int(),
		jen.Id("want").Any(),
		jen.Id("got").Any(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s call #%d argument \"%s\" did not match\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
				jen.Id("callNo"),
				jen.Id("argName"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("  want: %#v\n"),
				jen.Id("want"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("   got: %#v\n"),
				jen.Id("got"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("method: %s\n"),
				jen.Id("comparedBy"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(jen.Lit("\n")),
		jen.Id("m").Dot("buildCallHistory").Call(
			jen.Id("sb"),
			jen.Lit("call history"),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("hint: for custom matching use .Match(func(...) bool) at %s\n\tor use STUB for fine-grained control"),
				jen.Id("expectAt"),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageNotImplementedCode() jen.Code {
	return jen.Func().Id(lib.MessageNotImplementedFunc).Params(
		jen.Id("interfaceName").String(),
		jen.Id("methodName").String(),
		jen.Id("signature").String(),
		jen.Id("createdLocation").String(),
		jen.Id("args").Index().Any(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("unexpected call to %s.%s\n"),
				jen.Id("interfaceName"),
				jen.Id("methodName"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("signature: %s.%s%s\n"),
				jen.Id("interfaceName"),
				jen.Id("methodName"),
				jen.Id("signature"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("called at: %s\n"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(3)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(jen.Lit("arguments:\n")),
		jen.Id(lib.MessageWriteArgumentsFunc).Call(
			jen.Id("sb"),
			jen.Lit("\t%[MAX-KEY-LEN]s = %#v\n"),
			jen.Id("args"),
		),
		jen.Line(),
		jen.Id("location").Op(":=").Lit(""),
		jen.If(jen.Id("createdLocation").Op("!=").Lit("")).Block(
			jen.Id("location").Op("=").Lit(" after ").Op("+").Id("createdLocation"),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\nhint:%s use one of:\n\t[var].EXPECT().%s(t)\n\t[var].STUB().%s(func(...) ...)\n\n"),
				jen.Id("location"),
				jen.Id("methodName"),
				jen.Id("methodName"),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageCallHistoryCode() jen.Code {
	return jen.Func().Id(lib.MessageCallHistoryFunc).Params(
		jen.Id("sb").Op("*").Qual("strings", "Builder"),
		jen.Id("index").Int(),
		jen.Id("expectedAt").String(),
		jen.Id("calledAt").String(),
		jen.Id("args").Index().Any(),
	).String().Block(
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t#%d expect at: %s\n"),
				jen.Id("index").Op("+").Lit(1),
				jen.Id("expectedAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t   called at: %s\n"),
				jen.Id("calledAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t   arguments:\n"),
			),
		),
		jen.Id(lib.MessageWriteArgumentsFunc).Call(
			jen.Id("sb"),
			jen.Lit("\t\t%[MAX-KEY-LEN]s = %#v\n"),
			jen.Id("args"),
		),
		jen.Id("sb").Dot("WriteString").Call(jen.Lit("\n")),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageTooManyCallsCode() jen.Code {
	return jen.Func().Id(lib.MessageTooManyCallsFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("want").Int(),
		jen.Id("got").Int(),
		jen.Id("args").Index().Any(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("too many calls to %s.%s\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\twant: %d, got: %d\n\n"),
				jen.Id("want"),
				jen.Id("got"),
			),
		),
		jen.Id("m").Dot("buildCallHistory").Call(
			jen.Id("sb"),
			jen.Lit(""),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t#%d expect at: %s\n"),
				jen.Id("got"),
				jen.Lit("missing"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t   called at: %s\n"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(4)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t   arguments:\n"),
			),
		),
		jen.Id(lib.MessageWriteArgumentsFunc).Call(
			jen.Id("sb"),
			jen.Lit("\t\t%[MAX-KEY-LEN]s = %#v\n"),
			jen.Id("args"),
		),
		jen.Id("sb").Dot("WriteString").Call(jen.Lit("\n")),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\thint: remove unexpected call or add 1 more EXPECT:\n\t\t[var].EXPECT().%s(t)\n"),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageMatchByNilCode() jen.Code {
	return jen.Func().Id(lib.MessageMatchByNilFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s Match received a nil function\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("\thint: provide a valid function"),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageExpectByNilCode() jen.Code {
	return jen.Func().Id(lib.MessageExpectByNilFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("unexpected nil testing.TB in %s.%s\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\tcalled at: %s\n\n"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(3)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("\thint: EXPECT requires a valid testing.TB, use STUB instead:\n"),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t\tspy := [var].STUB().%s(func(...) ...)\n"),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageExpectAfterStubCode() jen.Code {
	return jen.Func().Id(lib.MessageExpectAfterStubFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("stubAt").String(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("conflicting usage for %s.%s\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n"),
				jen.Lit("STUB used at"),
				jen.Id("stubAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n\n"),
				jen.Lit("EXPECT used at"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(3)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("\thint: use either EXPECT or STUB for the same method, not both\n\n"),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageStubByNilCode() jen.Code {
	return jen.Func().Id(lib.MessageStubByNilFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("calledAt").String(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s STUB received a nil function\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("called at: %s\n\n"),
				jen.Id("calledAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("hint: provide a valid function\n"),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageStubAfterExpectCode() jen.Code {
	return jen.Func().Id(lib.MessageStubAfterExpectFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("expectAt").String(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("conflicting usage for %s.%s\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n"),
				jen.Lit("EXPECT used at"),
				jen.Id("expectAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n\n"),
				jen.Lit("STUB used at"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(3)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("\thint: use either EXPECT or STUB for the same method, not both\n\n"),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageDuplicateStubCode() jen.Code {
	return jen.Func().Id(lib.MessageDuplicateStubFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("firstUsedAt").String(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("duplicate STUB for %s.%s\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n"),
				jen.Lit("first used at"),
				jen.Id("firstUsedAt"),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t%14s: %s\n\n"),
				jen.Lit("second used at"),
				jen.Id(lib.CallerLocationFunc).Call(jen.Lit(3)),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\thint: %s.%s is already stubbed, remove one of the above\n\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) MessageExpectButNotCalledCode() jen.Code {
	return jen.Func().Id(lib.MessageExpectButNotCalledFunc).Params(
		jen.Id("m").Id(lib.MethodInterface),
		jen.Id("want").Int(),
		jen.Id("got").Int(),
		jen.Id("index").Int(),
	).String().Block(
		jen.Id("sb").Op(":=").Op("&").Qual("strings", "Builder").Values(),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s was not called as expected\n"),
				jen.Id("m").Dot("interfaceName").Call(),
				jen.Id("m").Dot("methodName").Call(),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\twant: %d, got: %d\n\n"),
				jen.Id("want"),
				jen.Id("got"),
			),
		),
		jen.Id("m").Dot("buildCallHistory").Call(
			jen.Id("sb"),
			jen.Lit(""),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("\t#%d never called\n\n"),
				jen.Id("index").Op("+").Lit(1),
			),
		),
		jen.Id("sb").Dot("WriteString").Call(
			jen.Lit("\thint: add the missing call or remove the EXPECT above"),
		),
		jen.Return(jen.Id("sb").Dot("String").Call()),
	)
}

func (lib *LibraryData) CompareByReflectEqualCode() jen.Code {
	return jen.Func().Id(lib.CompareByReflectEqualFunc).Types(
		jen.Id("M").Id(lib.MethodInterface),
		jen.Id("T").Any(),
	).Params(
		jen.Id("m").Id("M"),
		jen.Id("argName").String(),
		jen.Id("want").Id("T"),
		jen.Id("got").Id("T"),
		jen.Id("tb").Qual("testing", "TB"),
		jen.Id("expectAt").String(),
		jen.Id("index").Int(),
	).Block(
		jen.If(
			jen.Qual("reflect", "DeepEqual").Call(jen.Id("want"), jen.Id("got")),
		).Block(
			jen.Return(),
		),
		jen.Line(),
		jen.Id("tb").Dot("Helper").Call(),
		jen.Id("m").Dot("fatal").Call(
			jen.Id("index"),
			jen.Id(lib.MessageArgumentMismatchedFunc).Call(
				jen.Id("m"),
				jen.Id("argName"),
				jen.Id("expectAt"),
				jen.Lit("reflect.DeepEqual"),
				jen.Id("index").Op("+").Lit(1),
				jen.Id("want"),
				jen.Id("got"),
			),
		),
	)
}

func (lib *LibraryData) CompareByBasicComparisonCode() jen.Code {
	return jen.Func().Id(lib.CompareByBasicComparisonFunc).Types(
		jen.Id("M").Id(lib.MethodInterface),
		jen.Id("T").Comparable(),
	).Params(
		jen.Id("m").Id("M"),
		jen.Id("argName").String(),
		jen.Id("want").Id("T"),
		jen.Id("got").Id("T"),
		jen.Id("tb").Qual("testing", "TB"),
		jen.Id("expectAt").String(),
		jen.Id("index").Int(),
	).Block(
		jen.If(
			jen.Id("want").Op("==").Id("got"),
		).Block(
			jen.Return(),
		),
		jen.Line(),
		jen.Id("tb").Dot("Helper").Call(),
		jen.Id("m").Dot("fatal").Call(
			jen.Id("index"),
			jen.Id(lib.MessageArgumentMismatchedFunc).Call(
				jen.Id("m"),
				jen.Id("argName"),
				jen.Id("expectAt"),
				jen.Lit("=="),
				jen.Id("index").Op("+").Lit(1),
				jen.Id("want"),
				jen.Id("got"),
			),
		),
	)
}

func (lib *LibraryData) GenerateCode() []jen.Code {
	return []jen.Code{
		lib.CallerLocationCode(), jen.Line(), jen.Line(),
		lib.MethodInterfaceCode(), jen.Line(), jen.Line(),
		lib.MessageWriteArgumentsCode(), jen.Line(), jen.Line(),
		lib.MessageMatchFailCode(), jen.Line(), jen.Line(),
		lib.MessageArgumentMismatchedCode(), jen.Line(), jen.Line(),
		lib.MessageNotImplementedCode(), jen.Line(), jen.Line(),
		lib.MessageCallHistoryCode(), jen.Line(), jen.Line(),
		lib.MessageTooManyCallsCode(), jen.Line(), jen.Line(),
		lib.MessageMatchByNilCode(), jen.Line(), jen.Line(),
		lib.MessageExpectByNilCode(), jen.Line(), jen.Line(),
		lib.MessageExpectAfterStubCode(), jen.Line(), jen.Line(),
		lib.MessageStubByNilCode(), jen.Line(), jen.Line(),
		lib.MessageStubAfterExpectCode(), jen.Line(), jen.Line(),
		lib.MessageDuplicateStubCode(), jen.Line(), jen.Line(),
		lib.MessageExpectButNotCalledCode(), jen.Line(), jen.Line(),
		lib.CompareByReflectEqualCode(), jen.Line(), jen.Line(),
		lib.CompareByBasicComparisonCode(), jen.Line(), jen.Line(),
	}
}
