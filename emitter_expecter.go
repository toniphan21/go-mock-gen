package mockgen

//type targetExpecter struct { // skip:!expect
//	target *target
//}
//
//func (e *targetExpecter) Full(tb testing.TB) *targetFullExpecter { // skip:!expect
//	if e.target.td == nil {
//		e.target.td = &targetTestDouble{}
//	}
//
//	var mock = e.target.td.Full
//	if mock == nil {
//		mock = &targetFull{}
//		e.target.td.Full = mock
//	}
//
//	if mock.stub != nil {
//		mock.panic(libMessageExpectAfterStub(mock, mock.stubLocation))
//	}
//
//	if tb == nil {
//		mock.panic(libMessageExpectByNil(mock))
//	}
//
//	index := len(mock.expects)
//	mock.expects = append(mock.expects, &targetFullExpect{
//		location:         libCallerLocation(2),
//		matcher:          &targetFullArgumentMatcher{},
//		matcherWants:     make(map[string]any),
//		matcherMethods:   make(map[string]string),
//		matcherHints:     make(map[string]string),
//		matcherLocations: make(map[string]string),
//		index:            index,
//		tb:               tb,
//	})
//
//	tb.Helper()
//	tb.Cleanup(func() { tb.Helper(); mock.verify(index) })
//
//	return &targetFullExpecter{target: mock, expect: mock.expects[index]}
//}
