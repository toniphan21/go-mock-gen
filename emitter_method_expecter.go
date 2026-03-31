package mockgen

//type targetFullExpecter struct { // skip:!expect
//	target *targetFull
//	expect *targetFullExpect
//}
//
//func (e *targetFullExpecter) Return(first []Result, second error) { // skip:!expect
//	e.expect.returns = targetFullReturn{first: first, second: second}
//}
//
//func (e *targetFullExpecter) Match(matcher func(ctx context.Context, input string) bool) *targetFullExpecterWithMatch { // skip:!expect
//	if matcher == nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageMatchByNil(e.target))
//	}
//
//	e.expect.match = matcher
//	e.expect.matchLocation = libCallerLocation(2)
//	return &targetFullExpecterWithMatch{expect: e.expect}
//}
//
//func (e *targetFullExpecter) MatchCtx(matcher func(ctx context.Context) bool) *targetFullExpecterWithMatchArg { // skip:!expect
//	if matcher == nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchCtx"))
//	}
//
//	e.expect.matcher.ctx = matcher
//	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
//	e.expect.matcherHints["ctx"] = libMessageMatchArgHint()
//	return &targetFullExpecterWithMatchArg{expect: e.expect, target: e.target}
//}
//
//func (e *targetFullExpecter) MatchInput(matcher func(input string) bool) *targetFullExpecterWithMatchArg { // skip:!expect
//	if matcher == nil {
//		e.expect.tb.Helper()
//		e.target.fatal(e.expect.index, libMessageMatchArgByNil(e.target, "MatchInput"))
//	}
//
//	e.expect.matcher.input = matcher
//	e.expect.matcherLocations["input"] = libCallerLocation(2)
//	e.expect.matcherHints["input"] = libMessageMatchArgHint()
//	return &targetFullExpecterWithMatchArg{expect: e.expect, target: e.target}
//}
//
//func (e *targetFullExpecter) With(ctx context.Context, input string) *targetFullExpecterWithValue { // skip:!expect
//	e.WithCtx(ctx)
//	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
//
//	e.WithInput(input)
//	e.expect.matcherLocations["input"] = libCallerLocation(2)
//
//	return &targetFullExpecterWithValue{expect: e.expect}
//}
//
//func (e *targetFullExpecter) WithCtx(ctx context.Context) *targetFullExpecterWithValueArg { // skip:!expect
//	e.expect.matcher.ctx = libReflectEqualMatcher(ctx)
//	e.expect.matcherWants["ctx"] = ctx
//	e.expect.matcherMethods["ctx"] = "reflect.DeepEqual"
//	e.expect.matcherLocations["ctx"] = libCallerLocation(2)
//
//	return &targetFullExpecterWithValueArg{expect: e.expect, target: e.target}
//}
//
//func (e *targetFullExpecter) WithInput(input string) *targetFullExpecterWithValueArg { // skip:!expect
//	e.expect.matcher.input = libBasicComparisonMatcher(input)
//	e.expect.matcherWants["input"] = input
//	e.expect.matcherMethods["input"] = "=="
//	e.expect.matcherLocations["input"] = libCallerLocation(2)
//
//	return &targetFullExpecterWithValueArg{expect: e.expect, target: e.target}
//}
