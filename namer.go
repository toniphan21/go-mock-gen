package mockgen

import (
	"fmt"

	"github.com/ettle/strcase"
)

type Namer interface {
	Library() LibraryData
	Constructor() string
	Struct() string
	TestDouble() string
	Stubber() string
	Expecter() string
	Method(method string) MethodNamer
}

type MethodNamer interface {
	Struct() string
	Call() string
	Argument() string
	ArgumentField(originalName string, index int) (field string, name string)
	ArgumentMatcher() string
	Return() string
	ReturnField(originalName string, index int) (field string, name string)
	Expect() string
	Expecter() string
	ExpecterValue() string
	ExpecterMatch() string
	ExpecterValueArg() string
	ExpecterMatchArg() string
}

// ---

type NamerOption interface {
	apply(namer *defaultNamer)
}

type namerOptionFunc func(namer *defaultNamer)

func (f namerOptionFunc) apply(namer *defaultNamer) {
	f(namer)
}

func WithStructName(name string) NamerOption {
	return namerOptionFunc(func(namer *defaultNamer) {
		if name != "" {
			namer.structName = name
		}
	})
}

func WithLibraryPrefix(prefix string) NamerOption {
	return namerOptionFunc(func(namer *defaultNamer) {
		if prefix != "" {
			namer.libraryPrefix = prefix
		}
	})
}

func NewNamer(interfaceName string, options ...NamerOption) Namer {
	n := &defaultNamer{
		interfaceName: interfaceName,
		libraryPrefix: "mockgen",
		structName:    strcase.ToCamel(interfaceName),
	}

	for _, v := range options {
		v.apply(n)
	}

	return n
}

type defaultNamer struct {
	interfaceName string
	structName    string
	libraryPrefix string
}

func (n *defaultNamer) Library() LibraryData {
	return LibraryData{
		CallerLocationFunc:            n.libraryPrefix + "CallerLocation",
		MethodInterface:               n.libraryPrefix + "MockMethod",
		MessageWriteArgumentsFunc:     n.libraryPrefix + "MessageWriteArguments",
		MessageMatchFailFunc:          n.libraryPrefix + "MessageMatchFail",
		MessageNotImplementedFunc:     n.libraryPrefix + "MessageNotImplemented",
		MessageCallHistoryFunc:        n.libraryPrefix + "MessageCallHistory",
		MessageTooManyCallsFunc:       n.libraryPrefix + "MessageTooManyCalls",
		MessageMatchByNilFunc:         n.libraryPrefix + "MessageMatchByNil",
		MessageExpectByNilFunc:        n.libraryPrefix + "MessageExpectByNil",
		MessageExpectAfterStubFunc:    n.libraryPrefix + "MessageExpectAfterStub",
		MessageStubByNilFunc:          n.libraryPrefix + "MessageStubByNil",
		MessageStubAfterExpectFunc:    n.libraryPrefix + "MessageStubAfterExpect",
		MessageDuplicateStubFunc:      n.libraryPrefix + "MessageDuplicateStub",
		MessageExpectButNotCalledFunc: n.libraryPrefix + "MessageExpectButNotCalled",
		MessageMatchArgByNilFunc:      n.libraryPrefix + "MessageMatchArgByNil",
		MessageDuplicateMatchArgFunc:  n.libraryPrefix + "MessageDuplicateMatchArg",
		MessageMatchArgHintFunc:       n.libraryPrefix + "MessageMatchArgHint",
		MatchArgumentFunc:             n.libraryPrefix + "MatchArgument",
		ReflectEqualMatcherFunc:       n.libraryPrefix + "ReflectEqualMatcher",
		BasicComparisonMatcherFunc:    n.libraryPrefix + "BasicComparisonMatcher",
	}
}

func (n *defaultNamer) Constructor() string {
	return fmt.Sprintf("test%s", n.interfaceName)
}

func (n *defaultNamer) Struct() string {
	return n.structName
}

func (n *defaultNamer) TestDouble() string {
	return fmt.Sprintf("%sTestDouble", n.structName)
}

func (n *defaultNamer) Stubber() string {
	return fmt.Sprintf("%sStubber", n.structName)
}

func (n *defaultNamer) Expecter() string {
	return fmt.Sprintf("%sExpecter", n.structName)
}

func (n *defaultNamer) Method(method string) MethodNamer {
	return &defaultMethodNamer{base: fmt.Sprintf("%s%s", n.structName, method)}
}

var _ Namer = (*defaultNamer)(nil)

// ---

type defaultMethodNamer struct {
	base string
}

func (n *defaultMethodNamer) Struct() string {
	return n.base
}

func (n *defaultMethodNamer) Call() string {
	return fmt.Sprintf("%sCall", n.base)
}

func (n *defaultMethodNamer) Argument() string {
	return fmt.Sprintf("%sArgument", n.base)
}

func (n *defaultMethodNamer) ArgumentField(originalName string, index int) (field string, name string) {
	isExported := toPascalCase(n.base) == n.base
	if originalName != "" && originalName != "_" {
		field = originalName
		name = originalName
	} else {
		field = fmt.Sprintf("arg%d", index)
		name = field
	}

	if isExported {
		field = toPascalCase(field)
	}
	return field, name
}

func (n *defaultMethodNamer) ArgumentMatcher() string {
	return fmt.Sprintf("%sArgumentMatcher", n.base)
}

func (n *defaultMethodNamer) Return() string {
	return fmt.Sprintf("%sReturn", n.base)
}

func (n *defaultMethodNamer) ReturnField(originalName string, index int) (field string, name string) {
	ordinals := []string{"first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth", "tenth"}
	isExported := toPascalCase(n.base) == n.base
	if originalName != "" {
		field = originalName
		name = originalName
	} else {
		if index < len(ordinals) {
			field = ordinals[index]
			name = ordinals[index]
		} else {
			field = fmt.Sprintf("ret%d", index)
			name = field
		}
	}

	if isExported {
		field = toPascalCase(field)
	}
	return field, name
}

func (n *defaultMethodNamer) Expect() string {
	return fmt.Sprintf("%sExpect", n.base)
}

func (n *defaultMethodNamer) Expecter() string {
	return fmt.Sprintf("%sExpecter", n.base)
}

func (n *defaultMethodNamer) ExpecterValue() string {
	return fmt.Sprintf("%sExpecterWithValue", n.base)
}

func (n *defaultMethodNamer) ExpecterMatch() string {
	return fmt.Sprintf("%sExpecterWithMatch", n.base)
}

func (n *defaultMethodNamer) ExpecterValueArg() string {
	return fmt.Sprintf("%sExpecterWithValueArg", n.base)
}

func (n *defaultMethodNamer) ExpecterMatchArg() string {
	return fmt.Sprintf("%sExpecterWithMatchArg", n.base)
}

var _ MethodNamer = (*defaultMethodNamer)(nil)
