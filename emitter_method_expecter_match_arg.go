package mockgen

//type targetFullExpecterWithMatchArg struct { // skip:!expect
//	expect *targetFullExpect
//	target *targetFull
//}
//
//func (e *targetFullExpecterWithMatchArg) Return(first []Result, second error) { // skip:!expect
//	e.expect.returns = targetFullReturn{first: first, second: second}
//}
//
//func (e *targetFullExpecterWithMatchArg) MatchCtx(matcher func(ctx context.Context) bool) *targetFullExpecterWithMatchArg { // skip:!expect
//	if matcher == nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchCtx"))
//	}
//
//	if e.expect.matcher.ctx != nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchCtx", e.expect.matcherLocations["ctx"]))
//	}
//
//	e.expect.matcher.ctx = matcher
//	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
//	e.expect.matcherHints["ctx"] = libMessageMatchArgHint()
//	return e
//}
//
//func (e *targetFullExpecterWithMatchArg) MatchInput(matcher func(input string) bool) *targetFullExpecterWithMatchArg { // skip:!expect
//	if matcher == nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchInput"))
//	}
//
//	if e.expect.matcher.input != nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageDuplicateMatchArg(e.target, "MatchInput", e.expect.matcherLocations["input"]))
//	}
//
//	e.expect.matcher.input = matcher
//	e.expect.matcherLocations["input"] = libCallerLocation(2)
//	e.expect.matcherHints["input"] = libMessageMatchArgHint()
//	return e
//}
