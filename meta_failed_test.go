package mockgen

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhatp.com/go/mock-gen/internal/meta"
)

// Run this test first to generate the library output before running other tests.
// This acts as a regression check to ensure the generated library code works correctly.
func Test_GenerateCode_As_Regression_Test(t *testing.T) {
	lib := &LibraryData{
		CallerLocationFunc:            "libCallerLocation",
		MethodInterface:               "libMockMethod",
		MessageWriteArgumentsFunc:     "libMessageWriteArguments",
		MessageMatchFailFunc:          "libMessageMatchFail",
		MessageNotImplementedFunc:     "libMessageNotImplemented",
		MessageCallHistoryFunc:        "libMessageCallHistory",
		MessageTooManyCallsFunc:       "libMessageTooManyCalls",
		MessageMatchByNilFunc:         "libMessageMatchByNil",
		MessageExpectByNilFunc:        "libMessageExpectByNil",
		MessageExpectAfterStubFunc:    "libMessageExpectAfterStub",
		MessageStubByNilFunc:          "libMessageStubByNil",
		MessageStubAfterExpectFunc:    "libMessageStubAfterExpect",
		MessageDuplicateStubFunc:      "libMessageDuplicateStub",
		MessageExpectButNotCalledFunc: "libMessageExpectButNotCalled",
		MessageMatchArgByNilFunc:      "libMessageMatchArgByNil",
		MessageDuplicateMatchArgFunc:  "libMessageDuplicateMatchArg",
		MessageMatchArgHintFunc:       "libMessageMatchArgHint",
		MatchArgumentFunc:             "libMatchArgument",
		ReflectEqualMatcherFunc:       "libReflectEqualMatcher",
		BasicComparisonMatcherFunc:    "libBasicComparisonMatcher",
	}

	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	goldenPath := filepath.Join(testDir, "testdata", "meta", "regression_code_generated_by_LibraryData_emitter.go")

	jf := jen.NewFile("meta")
	codes := lib.GenerateCode()
	for _, code := range codes {
		jf.Add(code)
	}
	out := jf.GoString()

	err := os.WriteFile(goldenPath, []byte(out), 0644)
	require.NoError(t, err)
}

type metaFailedOutputTestCase struct {
	name     string
	test     string
	expected string
}

func (c *metaFailedOutputTestCase) Run(t *testing.T) {
	testDir := "./testdata/meta"
	testName := c.test
	if testName == "" {
		testName = "^" + c.name + "$"
	}

	cmd := exec.Command("go", "test", "-v", "-count=1", testDir, "-run", testName)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected the meta-test to fail, but it passed!")
	}

	assert.Equal(t, c.expected, meta.ExtractOutput(out))
}

func Test_MockFailureOutput_NotImplemented(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Not_Implemented",
			expected: `unexpected call to Target.Full
	signature: Target.Full(ctx Context, id string) ([]Result, error)
	called at: failed_not_impl_test.go:11
	arguments:
		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
		input = "anything"
	
	hint: after failed_not_impl_test.go:9 use one of:
		[var].EXPECT().Full(t)
		[var].STUB().Full(func(...) ...)`,
		},

		{
			name: "Test_Not_Implemented_Via_Struct",
			expected: `unexpected call to Target.Full
	signature: Target.Full(ctx Context, id string) ([]Result, error)
	called at: failed_not_impl_test.go:17
	arguments:
		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
		input = "anything"
	
	hint: use one of:
		[var].EXPECT().Full(t)
		[var].STUB().Full(func(...) ...)`,
		},

		{
			name: "Test_Not_Implemented_Via_Another_Func",
			expected: `unexpected call to Target.Full
	signature: Target.Full(ctx Context, id string) ([]Result, error)
	called at: failed_not_impl_test.go:27
	arguments:
		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
		input = "anything"
	
	hint: after failed_not_impl_test.go:21 use one of:
		[var].EXPECT().Full(t)
		[var].STUB().Full(func(...) ...)`,
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_BadUsage(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Use_STUB_Twice",
			expected: `duplicate STUB for Target.Full
		 first used at: failed_bad_usage_test.go:11
		second used at: failed_bad_usage_test.go:14
	
		hint: Target.Full is already stubbed, remove one of the above`,
		},

		{
			name: "Test_Use_STUB_Thrice",
			expected: `duplicate STUB for Target.Full
		 first used at: failed_bad_usage_test.go:22
		second used at: failed_bad_usage_test.go:26
	
		hint: Target.Full is already stubbed, remove one of the above`,
		},

		{
			name: "Test_Use_STUB_After_EXPECT",
			expected: `conflicting usage for Target.Full
		EXPECT used at: failed_bad_usage_test.go:38
		  STUB used at: failed_bad_usage_test.go:39
	
		hint: use either EXPECT or STUB for the same method, not both`,
		},

		{
			name: "Test_Use_EXPECT_After_STUB",
			expected: `conflicting usage for Target.Full
		  STUB used at: failed_bad_usage_test.go:47
		EXPECT used at: failed_bad_usage_test.go:50
	
		hint: use either EXPECT or STUB for the same method, not both`,
		},

		{
			name: "Test_Pass_Nil_To_EXPECT",
			expected: `unexpected nil testing.TB in Target.Full
		called at: failed_bad_usage_test.go:56
	
		hint: EXPECT requires a valid testing.TB, use STUB instead:
			spy := [var].STUB().Full(func(...) ...)`,
		},

		{
			name: "Test_Pass_Nil_To_STUB",
			expected: `Target.Full STUB received a nil function
	called at: failed_bad_usage_test.go:62
	
	hint: provide a valid function`,
		},

		{
			name: "Test_Pass_Nil_To_STUB_After_Expect",
			expected: `Target.Full STUB received a nil function
	called at: failed_bad_usage_test.go:69
	
	hint: provide a valid function`,
		},

		{
			name: "Test_Pass_Nil_To_EXPECT_Match",
			expected: `    failed_bad_usage_test.go:75: Target.Full Match received a nil function
        	hint: provide a valid function
`,
		},

		{
			name: "Test_Pass_Nil_To_EXPECT_After_EXPECT",
			expected: `unexpected nil testing.TB in Target.Full
		called at: failed_bad_usage_test.go:82
	
		hint: EXPECT requires a valid testing.TB, use STUB instead:
			spy := [var].STUB().Full(func(...) ...)`,
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_CalledMoreThanExpected(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_One_EXPECT_Call_Twice",
			expected: `too many calls to Target.Full
		want: 1, got: 2
	
		#1 expect at: failed_called_more_than_expected_test.go:11
		   called at: failed_called_more_than_expected_test.go:13
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "a"
	
		#2 expect at: missing
		   called at: failed_called_more_than_expected_test.go:14
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "b"
	
		hint: remove unexpected call or add 1 more EXPECT:
			[var].EXPECT().Full(t)`,
		},

		{
			name: "Test_Two_EXPECT_Call_Thrice",
			expected: `too many calls to Target.Full
		want: 2, got: 3
	
		#1 expect at: failed_called_more_than_expected_test.go:20
		   called at: failed_called_more_than_expected_test.go:23
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "a"
	
		#2 expect at: failed_called_more_than_expected_test.go:21
		   called at: failed_called_more_than_expected_test.go:24
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "b"
	
		#3 expect at: missing
		   called at: failed_called_more_than_expected_test.go:25
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "c"
	
		hint: remove unexpected call or add 1 more EXPECT:
			[var].EXPECT().Full(t)`,
		},

		{
			name: "Test_One_EXPECT_Call_Twice_In_Production",
			expected: `too many calls to Target.Full
		want: 1, got: 2
	
		#1 expect at: failed_called_more_than_expected_test.go:32
		   called at: production_code.go:18
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "any 1"
	
		#2 expect at: missing
		   called at: production_code.go:24
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "any 2"
	
		hint: remove unexpected call or add 1 more EXPECT:
			[var].EXPECT().Full(t)`,
		},

		{
			name: "Test_Two_EXPECT_Call_Thrice_In_Production",
			expected: `too many calls to Target.Full
		want: 2, got: 3
	
		#1 expect at: failed_called_more_than_expected_test.go:41
		   called at: production_code.go:28
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "any 1"
	
		#2 expect at: failed_called_more_than_expected_test.go:42
		   called at: production_code.go:34
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "any 2"
	
		#3 expect at: missing
		   called at: production_code.go:41
		   arguments:
			  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
			input = "any 3"
	
		hint: remove unexpected call or add 1 more EXPECT:
			[var].EXPECT().Full(t)`,
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_CalledLessThanExpected(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Two_EXPECT_Call_Once",
			expected: `    failed_call_less_than_expected_test.go:12: Target.Full was not called as expected
        	want: 2, got: 1
        
        	#1 expect at: failed_call_less_than_expected_test.go:11
        	   called at: failed_call_less_than_expected_test.go:14
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "a"
        
        	#2 never called
        
        	hint: add the missing call or remove the EXPECT above
`,
		},

		{
			name: "Test_Three_EXPECT_Call_Twice",
			expected: `    failed_call_less_than_expected_test.go:22: Target.Full was not called as expected
        	want: 3, got: 2
        
        	#1 expect at: failed_call_less_than_expected_test.go:20
        	   called at: failed_call_less_than_expected_test.go:24
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "a"
        
        	#2 expect at: failed_call_less_than_expected_test.go:21
        	   called at: failed_call_less_than_expected_test.go:25
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "b"
        
        	#3 never called
        
        	hint: add the missing call or remove the EXPECT above
`,
		},

		{
			name: "Test_Two_EXPECT_Call_Once_In_Production",
			expected: `    failed_call_less_than_expected_test.go:33: Target.Full was not called as expected
        	want: 2, got: 1
        
        	#1 expect at: failed_call_less_than_expected_test.go:32
        	   called at: production_code.go:14
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "any"
        
        	#2 never called
        
        	hint: add the missing call or remove the EXPECT above
`,
		},

		{
			name: "Test_Three_EXPECT_Call_Twice_In_Production",
			expected: `    failed_call_less_than_expected_test.go:44: Target.Full was not called as expected
        	want: 3, got: 2
        
        	#1 expect at: failed_call_less_than_expected_test.go:42
        	   called at: production_code.go:18
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "any 1"
        
        	#2 expect at: failed_call_less_than_expected_test.go:43
        	   called at: production_code.go:24
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "any 2"
        
        	#3 never called
        
        	hint: add the missing call or remove the EXPECT above
`,
		},

		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_Match(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Match_Fail_FirstCall",
			expected: `    failed_match_test.go:15: Target.Full call #1 did not match
        arguments:
        	  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        	input = "a"
        
        hint: check the callback passed to Match at failed_match_test.go:11
`,
		},

		{
			name: "Test_Match_Fail_SecondCall",
			expected: `    failed_match_test.go:29: Target.Full call #2 did not match
        arguments:
        	  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        	input = "a"
        
        call history:
        	#1 expect at: failed_match_test.go:21
        	   called at: failed_match_test.go:28
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "1"
        
        hint: check the callback passed to Match at failed_match_test.go:24
`,
		},

		{
			name: "Test_Match_Fail_FirstCall_Production",
			expected: `    production_code.go:14: Target.Full call #1 did not match
        arguments:
        	  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        	input = "a"
        
        hint: check the callback passed to Match at failed_match_test.go:36
`,
		},

		{
			name: "Test_Match_Fail_SecondCall_Production",
			expected: `    production_code.go:24: Target.Full call #2 did not match
        arguments:
        	  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        	input = "a 2"
        
        call history:
        	#1 expect at: failed_match_test.go:47
        	   called at: production_code.go:18
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "a 1"
        
        hint: check the callback passed to Match at failed_match_test.go:50
`,
		},
		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_CalledWith(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_CallWith_Fail_FirstCall_FirstArgument",
			expected: `    failed_called_with_test.go:14: Target.Full call #1 argument "ctx" did not match
          want: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
           got: &context.valueCtx{Context:context.backgroundCtx{emptyCtx:context.emptyCtx{}}, key:"key", val:"val"}
        method: reflect.DeepEqual
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:12
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_FirstCall_SecondArgument",
			expected: `    failed_called_with_test.go:23: Target.Full call #1 argument "input" did not match
          want: "a"
           got: "1"
        method: ==
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:21
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_SecondCall_FirstArgument",
			expected: `    failed_called_with_test.go:34: Target.Full call #2 argument "ctx" did not match
          want: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
           got: &context.valueCtx{Context:context.backgroundCtx{emptyCtx:context.emptyCtx{}}, key:"key", val:"val"}
        method: reflect.DeepEqual
        
        call history:
        	#1 expect at: failed_called_with_test.go:30
        	   called at: failed_called_with_test.go:33
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "1"
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:31
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_SecondCall_SecondArgument",
			expected: `    failed_called_with_test.go:45: Target.Full call #2 argument "input" did not match
          want: "1"
           got: "a"
        method: ==
        
        call history:
        	#1 expect at: failed_called_with_test.go:41
        	   called at: failed_called_with_test.go:44
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "1"
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:42
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_FirstCall_FirstArgument_Production",
			expected: `    production_code.go:14: Target.Full call #1 argument "ctx" did not match
          want: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
           got: &context.valueCtx{Context:context.backgroundCtx{emptyCtx:context.emptyCtx{}}, key:"key", val:"val"}
        method: reflect.DeepEqual
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:53
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_FirstCall_SecondArgument_Production",
			expected: `    production_code.go:14: Target.Full call #1 argument "input" did not match
          want: "a"
           got: "1"
        method: ==
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:63
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_SecondCall_FirstArgument_Production",
			expected: `    production_code.go:14: Target.Full call #2 argument "ctx" did not match
          want: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
           got: &context.valueCtx{Context:context.backgroundCtx{emptyCtx:context.emptyCtx{}}, key:"key", val:"val"}
        method: reflect.DeepEqual
        
        call history:
        	#1 expect at: failed_called_with_test.go:73
        	   called at: failed_called_with_test.go:76
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "1"
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:74
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_CallWith_Fail_SecondCall_SecondArgument_Production",
			expected: `    production_code.go:14: Target.Full call #2 argument "input" did not match
          want: "1"
           got: "a"
        method: ==
        
        call history:
        	#1 expect at: failed_called_with_test.go:85
        	   called at: failed_called_with_test.go:88
        	   arguments:
        		  ctx = context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        		input = "1"
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_called_with_test.go:86
        	or use STUB for fine-grained control
`,
		},
		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_MockFailureOutput_Stub(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Stub_Fail_Not_Called",
			expected: `    failed_stub_test.go:16: want 1, got 0
`,
		},

		{
			name: "Test_Stub_Fail_Too_Many_Calls",
			expected: `    failed_stub_test.go:31: want 1, got 2
`,
		},
		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_Arg_Matcher(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Arg_Matcher_Pass_Nil_In_First_Place_FirstArg",
			expected: `    failed_arg_matcher_test.go:11: Target.Full MatchCtx received a nil function
        	hint: provide a valid function
`,
		},

		{
			name: "Test_Arg_Matcher_Pass_Nil_In_Second_Place_FirstArg",
			expected: `    failed_arg_matcher_test.go:19: Target.Full MatchCtx received a nil function
        	hint: provide a valid function
`,
		},

		{
			name: "Test_Arg_Matcher_Pass_Nil_In_First_Place_SecondArg",
			expected: `    failed_arg_matcher_test.go:25: Target.Full MatchInput received a nil function
        	hint: provide a valid function
`,
		},

		{
			name: "Test_Arg_Matcher_Pass_Nil_In_Second_Place_SecondArg",
			expected: `    failed_arg_matcher_test.go:33: Target.Full MatchInput received a nil function
        	hint: provide a valid function
`,
		},

		{
			name: "Test_Arg_Matcher_Duplicate_FirstArg",
			expected: `    failed_arg_matcher_test.go:41: duplicate MatchCtx for Target.Full
        	 first used at: failed_arg_matcher_test.go:40
        
        	hint: each argument can only be matched once, remove one of the above
`,
		},

		{
			name: "Test_Arg_Matcher_Duplicate_SecondArg",
			expected: `    failed_arg_matcher_test.go:52: duplicate MatchInput for Target.Full
        	 first used at: failed_arg_matcher_test.go:51
        
        	hint: each argument can only be matched once, remove one of the above
`,
		},

		{
			name: "Test_Arg_Matcher_Failed_FirstPlace_FirstArg",
			expected: `    failed_arg_matcher_test.go:63: Target.Full call #1 argument "ctx" did not match
           got: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        method: func(got) bool
        
        	hint: check argument matching at failed_arg_matcher_test.go:61
        		or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Matcher_Failed_SecondPlace_FirstArg",
			expected: `    failed_arg_matcher_test.go:73: Target.Full call #1 argument "ctx" did not match
           got: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
        method: func(got) bool
        
        	hint: check argument matching at failed_arg_matcher_test.go:71
        		or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Matcher_Failed_FirstPlace_SecondArg",
			expected: `    failed_arg_matcher_test.go:82: Target.Full call #1 argument "input" did not match
           got: "anything"
        method: func(got) bool
        
        	hint: check argument matching at failed_arg_matcher_test.go:80
        		or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Matcher_Failed_SecondPlace_SecondArg",
			expected: `    failed_arg_matcher_test.go:92: Target.Full call #1 argument "input" did not match
           got: "anything"
        method: func(got) bool
        
        	hint: check argument matching at failed_arg_matcher_test.go:90
        		or use STUB for fine-grained control
`,
		},
		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

func Test_Arg_Value(t *testing.T) {
	cases := []metaFailedOutputTestCase{
		{
			name: "Test_Arg_Value_Duplicate_FirstArg",
			expected: `    failed_arg_value_test.go:13: duplicate WithCtx for Target.Full
        	 first used at: failed_arg_value_test.go:12
        
        	hint: each argument can only be matched once, remove one of the above
`,
		},

		{
			name: "Test_Arg_Value_Duplicate_SecondArg",
			expected: `    failed_arg_value_test.go:24: duplicate WithInput for Target.Full
        	 first used at: failed_arg_value_test.go:23
        
        	hint: each argument can only be matched once, remove one of the above
`,
		},

		{
			name: "Test_Arg_Value_Failed_FirstPlace_FirstArg",
			expected: `    failed_arg_value_test.go:35: Target.Full call #1 argument "ctx" did not match
          want: context.backgroundCtx{emptyCtx:context.emptyCtx{}}
           got: &context.valueCtx{Context:context.backgroundCtx{emptyCtx:context.emptyCtx{}}, key:"key", val:"val"}
        method: reflect.DeepEqual
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_arg_value_test.go:33
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Value_Failed_SecondPlace_FirstArg",
			expected: `    failed_arg_value_test.go:45: Target.Full call #1 argument "input" did not match
          want: "a"
           got: "anything"
        method: ==
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_arg_value_test.go:43
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Value_Failed_FirstPlace_SecondArg",
			expected: `    failed_arg_value_test.go:54: Target.Full call #1 argument "input" did not match
          want: "a"
           got: "anything"
        method: ==
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_arg_value_test.go:52
        	or use STUB for fine-grained control
`,
		},

		{
			name: "Test_Arg_Value_Failed_SecondPlace_SecondArg",
			expected: `    failed_arg_value_test.go:64: Target.Full call #1 argument "input" did not match
          want: "a"
           got: "anything"
        method: ==
        
        hint: for custom matching use .Match[arg](func(...) bool) at failed_arg_value_test.go:62
        	or use STUB for fine-grained control
`,
		},
		// ---
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.Run(t)
		})
	}
}
