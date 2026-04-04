# go-mock-gen

`go-mock-gen` is a tool to generate a type-safe, minimal, zero-dependency and idiomatic mock for testing with focus on
the Developer Experience. 

No `any`, no magic matchers. The API is designed so you can't do anything wrong - and when you do, it tells you exactly
why, where, and how to fix it.

image quick look

## Quick Usage

~~~bash
go run nhatp.com/go/mock-gen/cmd/go-mock-gen -i Interface -o mock_interface_test.go
~~~

---

## API

There are 2 API categories:

- The `.EXPECT()` way is for convenience.
- The `.STUB()` way is for fine-grain control.

On same method you cannot mix between them otherwise the test will fail immediately.

### EXPECT

There are only 10 ways to set an expectation - no Once(), no Twice(), no Times(). If you want to expect 2 calls, just
expect twice.

| after `.EXPECT().Method(t)`           | Arguments         | Return | Usage                        |
|---------------------------------------|-------------------|--------|------------------------------|
| `<empty>`                             | -                 | zero   | expect the call, ignore args |
| `.Return(…)`                          | -                 | …      | ignore args, stub return     |
| `.With(…)`                            | all, value        | zero   | match all args by value      |
| `.With(…).Return(…)`                  | all, value        | …      | match all args by value      |
| `.With[Arg](…)`                       | partial, value    | zero   | match ctx by value           |
| `.With[Arg](…).Return(…)`             | partial, value    | …      | match single arg by value    |
| `.Match(func(…) bool)`                | all, callback     | zero   | match all args by callback   |
| `.Match(func(…) bool).Return(…)`      | all, callback     | …      | match all args by callback   |
| `.Match[Arg](func(…) bool)`           | partial, callback | zero   | match single arg by callback |
| `.Match[Arg](func(…) bool).Return(…)` | partial, callback | …      | match single arg by callback |

If you use it in a wrong way the IDE will show you the error. In case it is not a syntax error the test will fail and
show you exactly why.

```go
package test

import "testing"

func Test_Quick_Expect_Example(t *testing.T) {
	repo := testRepository()

	t.Run("expect one call - ignore args - return zero", func(t *testing.T) {
		repo.EXPECT().GetUsers(t)
		// ...
		repo.GetUsers(...)
	})

	t.Run("expect two calls - first call match arg - second call stub return", func(t *testing.T) {
		repo.EXPECT().GetUsers(t).With(...)
		repo.EXPECT().GetUsers(t).Return(...)
		// ...
		repo.GetUsers(...)
		out := repo.GetUsers(...)
	})
}

```

### STUB

Stub API even simpler than the EXPECT, you need to provide the STUB function which has the same signature with the
implementation, it returns you a spy for you to assert yourself:

```go
package test

import (
	"fmt"
	"testing"
)

func Test_STUB(t *testing.T) {
	repo := testRepository()

	spy := repo.STUB().GetUsers(func(input string) []string {
		// inspect calls yourself
		if input != "awesome" {
			t.Fatal("hey, be awesome")
		}
		return nil
	})

	// ...test code...

	// inspect calls yourself
	if len(spy.Calls) != 1 {
		t.Fatal("why didn't you call me?")
	}

	if spy.Calls[0].Arguments.Input == "awesome" {
		fmt.Println("good!")
	}
}
```

---

## DX focus

