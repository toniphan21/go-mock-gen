package mockgen

import (
	"embed"
	"io/fs"
	"path"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/emitter/lib.*.go
var goldenEmitterFiles embed.FS

func Test_LibraryData(t *testing.T) {
	cases := []struct {
		file string
		data LibraryData
		fn   func(data LibraryData) jen.Code
	}{
		{
			file: "lib.CallerLocationCode.go",
			data: LibraryData{
				CallerLocationFunc: "repositoryCallerLocation",
			},
			fn: func(data LibraryData) jen.Code {
				return data.CallerLocationCode()
			},
		},

		{
			file: "lib.MethodInterfaceCode.go",
			data: LibraryData{
				MethodInterface: "repositoryMockMethod",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MethodInterfaceCode()
			},
		},

		{
			file: "lib.MessageWriteArgumentsCode.go",
			data: LibraryData{
				MessageWriteArgumentsFunc: "repositoryMessageWriteArguments",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageWriteArgumentsCode()
			},
		},

		{
			file: "lib.MessageMatchFailCode.go",
			data: LibraryData{
				MethodInterface:           "repositoryMockMethod",
				MessageWriteArgumentsFunc: "repositoryMessageWriteArguments",
				MessageMatchFailFunc:      "repositoryMessageMatchFail",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageMatchFailCode()
			},
		},

		{
			file: "lib.MessageArgumentMismatchedCode.go",
			data: LibraryData{
				MethodInterface:               "repositoryMockMethod",
				CallerLocationFunc:            "repositoryCallerLocation",
				MessageArgumentMismatchedFunc: "repositoryMessageArgumentMismatched",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageArgumentMismatchedCode()
			},
		},

		{
			file: "lib.MessageNotImplementedCode.go",
			data: LibraryData{
				CallerLocationFunc:        "repositoryCallerLocation",
				MessageNotImplementedFunc: "repositoryMessageNotImplemented",
				MessageWriteArgumentsFunc: "repositoryMessageWriteArguments",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageNotImplementedCode()
			},
		},

		{
			file: "lib.MessageCallHistoryCode.go",
			data: LibraryData{
				MessageCallHistoryFunc:    "repositoryMessageCallHistory",
				MessageWriteArgumentsFunc: "repositoryMessageWriteArguments",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageCallHistoryCode()
			},
		},

		{
			file: "lib.MessageTooManyCallsCode.go",
			data: LibraryData{
				MethodInterface:           "repositoryMockMethod",
				CallerLocationFunc:        "repositoryCallerLocation",
				MessageTooManyCallsFunc:   "repositoryMessageTooManyCalls",
				MessageWriteArgumentsFunc: "repositoryMessageWriteArguments",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageTooManyCallsCode()
			},
		},

		{
			file: "lib.MessageMatchByNilCode.go",
			data: LibraryData{
				MethodInterface:       "repositoryMockMethod",
				MessageMatchByNilFunc: "repositoryMessageMatchByNil",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageMatchByNilCode()
			},
		},

		{
			file: "lib.MessageExpectByNilCode.go",
			data: LibraryData{
				MethodInterface:        "repositoryMockMethod",
				CallerLocationFunc:     "repositoryCallerLocation",
				MessageExpectByNilFunc: "repositoryMessageExpectByNil",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageExpectByNilCode()
			},
		},

		{
			file: "lib.MessageExpectAfterStubCode.go",
			data: LibraryData{
				MethodInterface:            "repositoryMockMethod",
				CallerLocationFunc:         "repositoryCallerLocation",
				MessageExpectAfterStubFunc: "repositoryMessageExpectAfterStub",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageExpectAfterStubCode()
			},
		},

		{
			file: "lib.MessageStubByNilCode.go",
			data: LibraryData{
				MethodInterface:      "repositoryMockMethod",
				MessageStubByNilFunc: "repositoryMessageStubByNil",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageStubByNilCode()
			},
		},

		{
			file: "lib.MessageStubAfterExpectCode.go",
			data: LibraryData{
				MethodInterface:            "repositoryMockMethod",
				CallerLocationFunc:         "repositoryCallerLocation",
				MessageStubAfterExpectFunc: "repositoryMessageStubAfterExpect",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageStubAfterExpectCode()
			},
		},

		{
			file: "lib.MessageDuplicateStubCode.go",
			data: LibraryData{
				MethodInterface:          "repositoryMockMethod",
				CallerLocationFunc:       "repositoryCallerLocation",
				MessageDuplicateStubFunc: "repositoryMessageDuplicateStub",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageDuplicateStubCode()
			},
		},

		{
			file: "lib.MessageExpectButNotCalledCode.go",
			data: LibraryData{
				MethodInterface:               "repositoryMockMethod",
				MessageExpectButNotCalledFunc: "repositoryMessageExpectButNotCalled",
			},
			fn: func(data LibraryData) jen.Code {
				return data.MessageExpectButNotCalledCode()
			},
		},

		{
			file: "lib.CompareByReflectEqualCode.go",
			data: LibraryData{
				MethodInterface:               "repositoryMockMethod",
				MessageArgumentMismatchedFunc: "repositoryMessageArgumentMismatched",
				CompareByReflectEqualFunc:     "repositoryCompareByReflectEqual",
			},
			fn: func(data LibraryData) jen.Code {
				return data.CompareByReflectEqualCode()
			},
		},

		{
			file: "lib.CompareByBasicComparisonCode.go",
			data: LibraryData{
				MethodInterface:               "repositoryMockMethod",
				MessageArgumentMismatchedFunc: "repositoryMessageArgumentMismatched",
				CompareByBasicComparisonFunc:  "repositoryCompareByBasicComparison",
			},
			fn: func(data LibraryData) jen.Code {
				return data.CompareByBasicComparisonCode()
			},
		},

		{
			file: "lib.service.go",
			data: LibraryData{
				CallerLocationFunc:            "serviceCallerLocation",
				MethodInterface:               "serviceMockMethod",
				MessageWriteArgumentsFunc:     "serviceMessageWriteArguments",
				MessageMatchFailFunc:          "serviceMessageMatchFail",
				MessageArgumentMismatchedFunc: "serviceMessageArgumentMismatched",
				MessageNotImplementedFunc:     "serviceMessageNotImplemented",
				MessageCallHistoryFunc:        "serviceMessageCallHistory",
				MessageTooManyCallsFunc:       "serviceMessageTooManyCalls",
				MessageMatchByNilFunc:         "serviceMessageMatchByNil",
				MessageExpectByNilFunc:        "serviceMessageExpectByNil",
				MessageExpectAfterStubFunc:    "serviceMessageExpectAfterStub",
				MessageStubByNilFunc:          "serviceMessageStubByNil",
				MessageStubAfterExpectFunc:    "serviceMessageStubAfterExpect",
				MessageDuplicateStubFunc:      "serviceMessageDuplicateStub",
				MessageExpectButNotCalledFunc: "serviceMessageExpectButNotCalled",
				CompareByReflectEqualFunc:     "serviceCompareByReflectEqual",
				CompareByBasicComparisonFunc:  "serviceCompareByBasicComparison",
			},
			fn: func(data LibraryData) jen.Code {
				return jen.Add(data.GenerateCode()...)
			},
		},
		// ---
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			jf := jen.NewFile("emitter")
			code := tc.fn(tc.data)
			if code != nil {
				jf.Add(code)
			}
			out := jf.GoString()

			content, err := fs.ReadFile(goldenEmitterFiles, path.Join("testdata", "emitter", tc.file))
			require.NoError(t, err)

			assert.Equal(t, string(content), out)
		})
	}
}
