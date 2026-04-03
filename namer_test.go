package mockgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewNamer(t *testing.T) {
	t.Run("without any options", func(t *testing.T) {
		n, ok := NewNamer("Repository").(*defaultNamer)
		require.True(t, ok)

		assert.Equal(t, "Repository", n.interfaceName)
		assert.Equal(t, "repository", n.structName)
		assert.Equal(t, "mockgen", n.libraryPrefix)
	})

	t.Run("with struct name option", func(t *testing.T) {
		n, ok := NewNamer("Repository", WithStructName("")).(*defaultNamer)
		require.True(t, ok)
		assert.Equal(t, "repository", n.structName)

		n, ok = NewNamer("Repository", WithStructName("Repository")).(*defaultNamer)
		require.True(t, ok)
		assert.Equal(t, "Repository", n.structName)
	})

	t.Run("with library prefix option", func(t *testing.T) {
		n, ok := NewNamer("Repository", WithLibraryPrefix("")).(*defaultNamer)
		require.True(t, ok)
		assert.Equal(t, "mockgen", n.libraryPrefix)

		n, ok = NewNamer("Repository", WithLibraryPrefix("repository")).(*defaultNamer)
		require.True(t, ok)
		assert.Equal(t, "repository", n.libraryPrefix)
	})
}

func Test_DefaultNamer(t *testing.T) {
	n := NewNamer("Repository")

	assert.Equal(t, LibraryData{
		CallerLocationFunc:            "mockgenCallerLocation",
		MethodInterface:               "mockgenMockMethod",
		MessageWriteArgumentsFunc:     "mockgenMessageWriteArguments",
		MessageMatchFailFunc:          "mockgenMessageMatchFail",
		MessageNotImplementedFunc:     "mockgenMessageNotImplemented",
		MessageCallHistoryFunc:        "mockgenMessageCallHistory",
		MessageTooManyCallsFunc:       "mockgenMessageTooManyCalls",
		MessageMatchByNilFunc:         "mockgenMessageMatchByNil",
		MessageExpectByNilFunc:        "mockgenMessageExpectByNil",
		MessageExpectAfterStubFunc:    "mockgenMessageExpectAfterStub",
		MessageStubByNilFunc:          "mockgenMessageStubByNil",
		MessageStubAfterExpectFunc:    "mockgenMessageStubAfterExpect",
		MessageDuplicateStubFunc:      "mockgenMessageDuplicateStub",
		MessageExpectButNotCalledFunc: "mockgenMessageExpectButNotCalled",
		MessageMatchArgByNilFunc:      "mockgenMessageMatchArgByNil",
		MessageDuplicateMatchArgFunc:  "mockgenMessageDuplicateMatchArg",
		MessageMatchArgHintFunc:       "mockgenMessageMatchArgHint",
		MatchArgumentFunc:             "mockgenMatchArgument",
		ReflectEqualMatcherFunc:       "mockgenReflectEqualMatcher",
		BasicComparisonMatcherFunc:    "mockgenBasicComparisonMatcher",
	}, n.Library())
	assert.Equal(t, "testRepository", n.Constructor())
	assert.Equal(t, "repository", n.Struct())
	assert.Equal(t, "repositoryTestDouble", n.TestDouble())
	assert.Equal(t, "repositoryStubber", n.Stubber())
	assert.Equal(t, "repositoryExpecter", n.Expecter())

	m, ok := n.Method("GetUsers").(*defaultMethodNamer)
	require.True(t, ok)

	assert.Equal(t, "repositoryGetUsers", m.base)
}

func Test_DefaultMethodNamer(t *testing.T) {
	t.Run("named funcs", func(t *testing.T) {
		n := NewNamer("Repository")
		m := n.Method("GetUsers")

		assert.Equal(t, "repositoryGetUsers", m.Struct())
		assert.Equal(t, "repositoryGetUsersCall", m.Call())
		assert.Equal(t, "repositoryGetUsersArgument", m.Argument())
		assert.Equal(t, "repositoryGetUsersArgumentMatcher", m.ArgumentMatcher())
		assert.Equal(t, "repositoryGetUsersReturn", m.Return())
		assert.Equal(t, "repositoryGetUsersExpect", m.Expect())
		assert.Equal(t, "repositoryGetUsersExpecter", m.Expecter())
		assert.Equal(t, "repositoryGetUsersExpecterWithValue", m.ExpecterValue())
		assert.Equal(t, "repositoryGetUsersExpecterWithValueArg", m.ExpecterValueArg())
		assert.Equal(t, "repositoryGetUsersExpecterWithMatch", m.ExpecterMatch())
		assert.Equal(t, "repositoryGetUsersExpecterWithMatchArg", m.ExpecterMatchArg())
	})

	t.Run("ArgumentField", func(t *testing.T) {
		cases := []struct {
			base         string
			originalName string
			index        int
			field        string
			name         string
		}{
			{base: "repository", index: 0, field: "arg0", name: "arg0"},
			{base: "repository", index: 1, field: "arg1", name: "arg1"},
			{base: "repository", index: 2, field: "arg2", name: "arg2"},
			{base: "repository", index: 3, field: "arg3", name: "arg3"},
			{base: "repository", index: 4, field: "arg4", name: "arg4"},

			{base: "Repository", index: 0, field: "Arg0", name: "arg0"},
			{base: "Repository", index: 1, field: "Arg1", name: "arg1"},
			{base: "Repository", index: 2, field: "Arg2", name: "arg2"},
			{base: "Repository", index: 3, field: "Arg3", name: "arg3"},
			{base: "Repository", index: 4, field: "Arg4", name: "arg4"},

			{base: "repository", originalName: "_", index: 0, field: "arg0", name: "arg0"},
			{base: "repository", originalName: "_", index: 1, field: "arg1", name: "arg1"},
			{base: "repository", originalName: "_", index: 2, field: "arg2", name: "arg2"},
			{base: "repository", originalName: "_", index: 3, field: "arg3", name: "arg3"},
			{base: "repository", originalName: "_", index: 4, field: "arg4", name: "arg4"},

			{base: "Repository", originalName: "_", index: 0, field: "Arg0", name: "arg0"},
			{base: "Repository", originalName: "_", index: 1, field: "Arg1", name: "arg1"},
			{base: "Repository", originalName: "_", index: 2, field: "Arg2", name: "arg2"},
			{base: "Repository", originalName: "_", index: 3, field: "Arg3", name: "arg3"},
			{base: "Repository", originalName: "_", index: 4, field: "Arg4", name: "arg4"},

			{base: "repository", originalName: "var0", index: 0, field: "var0", name: "var0"},
			{base: "repository", originalName: "var1", index: 1, field: "var1", name: "var1"},
			{base: "repository", originalName: "var2", index: 2, field: "var2", name: "var2"},
			{base: "repository", originalName: "var3", index: 3, field: "var3", name: "var3"},
			{base: "repository", originalName: "var4", index: 4, field: "var4", name: "var4"},

			{base: "repository", originalName: "Var0", index: 0, field: "Var0", name: "Var0"},
			{base: "repository", originalName: "Var1", index: 1, field: "Var1", name: "Var1"},
			{base: "repository", originalName: "Var2", index: 2, field: "Var2", name: "Var2"},
			{base: "repository", originalName: "Var3", index: 3, field: "Var3", name: "Var3"},
			{base: "repository", originalName: "Var4", index: 4, field: "Var4", name: "Var4"},

			{base: "Repository", originalName: "var0", index: 0, field: "Var0", name: "var0"},
			{base: "Repository", originalName: "var1", index: 1, field: "Var1", name: "var1"},
			{base: "Repository", originalName: "var2", index: 2, field: "Var2", name: "var2"},
			{base: "Repository", originalName: "var3", index: 3, field: "Var3", name: "var3"},
			{base: "Repository", originalName: "var4", index: 4, field: "Var4", name: "var4"},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%s - originalName: %#v - index %d", tc.base, tc.originalName, tc.index), func(t *testing.T) {
				m := &defaultMethodNamer{base: tc.base}
				field, name := m.ArgumentField(tc.originalName, tc.index)

				assert.Equal(t, tc.field, field)
				assert.Equal(t, tc.name, name)
			})
		}
	})

	t.Run("ReturnField", func(t *testing.T) {
		cases := []struct {
			base         string
			originalName string
			index        int
			field        string
			name         string
		}{
			{base: "repository", index: 0, field: "first", name: "first"},
			{base: "repository", index: 1, field: "second", name: "second"},
			{base: "repository", index: 2, field: "third", name: "third"},
			{base: "repository", index: 3, field: "fourth", name: "fourth"},
			{base: "repository", index: 4, field: "fifth", name: "fifth"},
			{base: "repository", index: 5, field: "sixth", name: "sixth"},
			{base: "repository", index: 6, field: "seventh", name: "seventh"},
			{base: "repository", index: 7, field: "eighth", name: "eighth"},
			{base: "repository", index: 8, field: "ninth", name: "ninth"},
			{base: "repository", index: 9, field: "tenth", name: "tenth"},
			{base: "repository", index: 10, field: "ret10", name: "ret10"},
			{base: "repository", index: 11, field: "ret11", name: "ret11"},

			{base: "Repository", index: 0, field: "First", name: "first"},
			{base: "Repository", index: 1, field: "Second", name: "second"},
			{base: "Repository", index: 2, field: "Third", name: "third"},
			{base: "Repository", index: 3, field: "Fourth", name: "fourth"},
			{base: "Repository", index: 4, field: "Fifth", name: "fifth"},
			{base: "Repository", index: 5, field: "Sixth", name: "sixth"},
			{base: "Repository", index: 6, field: "Seventh", name: "seventh"},
			{base: "Repository", index: 7, field: "Eighth", name: "eighth"},
			{base: "Repository", index: 8, field: "Ninth", name: "ninth"},
			{base: "Repository", index: 9, field: "Tenth", name: "tenth"},
			{base: "Repository", index: 10, field: "Ret10", name: "ret10"},
			{base: "Repository", index: 11, field: "Ret11", name: "ret11"},

			{base: "repository", originalName: "ret0", index: 0, field: "ret0", name: "ret0"},
			{base: "repository", originalName: "ret1", index: 1, field: "ret1", name: "ret1"},
			{base: "repository", originalName: "ret2", index: 2, field: "ret2", name: "ret2"},
			{base: "repository", originalName: "ret3", index: 3, field: "ret3", name: "ret3"},
			{base: "repository", originalName: "ret4", index: 4, field: "ret4", name: "ret4"},
			{base: "repository", originalName: "ret5", index: 5, field: "ret5", name: "ret5"},

			{base: "repository", originalName: "Ret0", index: 0, field: "Ret0", name: "Ret0"},
			{base: "repository", originalName: "Ret1", index: 1, field: "Ret1", name: "Ret1"},
			{base: "repository", originalName: "Ret2", index: 2, field: "Ret2", name: "Ret2"},
			{base: "repository", originalName: "Ret3", index: 3, field: "Ret3", name: "Ret3"},
			{base: "repository", originalName: "Ret4", index: 4, field: "Ret4", name: "Ret4"},
			{base: "repository", originalName: "Ret5", index: 5, field: "Ret5", name: "Ret5"},

			{base: "Repository", originalName: "ret0", index: 0, field: "Ret0", name: "ret0"},
			{base: "Repository", originalName: "ret1", index: 1, field: "Ret1", name: "ret1"},
			{base: "Repository", originalName: "ret2", index: 2, field: "Ret2", name: "ret2"},
			{base: "Repository", originalName: "ret3", index: 3, field: "Ret3", name: "ret3"},
			{base: "Repository", originalName: "ret4", index: 4, field: "Ret4", name: "ret4"},
			{base: "Repository", originalName: "ret5", index: 5, field: "Ret5", name: "ret5"},
		}

		for _, tc := range cases {
			t.Run(fmt.Sprintf("%s - originalName: %#v - index %d", tc.base, tc.originalName, tc.index), func(t *testing.T) {
				m := &defaultMethodNamer{base: tc.base}
				field, name := m.ReturnField(tc.originalName, tc.index)

				assert.Equal(t, tc.field, field)
				assert.Equal(t, tc.name, name)
			})
		}
	})
}
