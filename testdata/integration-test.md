## Integration test

set up a golang project

```go.mod
module github.com/you/project

go 1.24
```

given the source code

```go
// file: source.go

package domain

import "context"

type Service interface {
	CreateUser(ctx context.Context, username string) error
}
```

run this command
```bash
// file: generate.sh
#!/bin/sh

go run nhatp.com/go/mock-gen/cmd/go-mock-gen \
    --interface Service -o output.go
```

it will generate this file, but we do not care because generated content is not important, this is an
integration test which verifies real output behavior.

```go
// golden-file: output.go
// ignored-content

package domain
```

given two tests

```go
// file: source_test.go

package domain

import (
	"context"
	"testing"
)

func Test_Service_CreateUser_Fail(t *testing.T) {
	mock := testService()
	mock.EXPECT().CreateUser(t)

	// do not call to make the test failed
}

func Test_Service_CreateUser_Pass(t *testing.T) {
	mock := testService()
	mock.EXPECT().CreateUser(t)

	mock.CreateUser(context.Background(), "username")
}
```

expect the tests results
```txt
// file: expected_test_result.txt
=== RUN  Test_Service_CreateUser_Fail
    source_test.go:11: Service.CreateUser was not called as expected
        	want: 1, got: 0
        
        	#1 never called
        
        	hint: add the missing call or remove the EXPECT above
--- FAIL

=== RUN  Test_Service_CreateUser_Pass
--- PASS
```