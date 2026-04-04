package mockgen

import (
	"testing"

	"nhatp.com/go/gen-lib/gentest"
)

func Test_Example_GenerateCode(t *testing.T) {
	cases := []struct {
		name     string
		data     ExampleData
		expected string
	}{
		{
			name: "no arguments, no returns",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				SkipExpect:    false,
			},
			expected: `package emitter

import (
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock := testTarget()

	t.Run("expect called once", func(t *testing.T) {
		mock.EXPECT().Method(t)

		mock.Method()
	})

	t.Run("expect called twice", func(t *testing.T) {
		mock.EXPECT().Method(t)
		mock.EXPECT().Method(t)

		mock.Method()
		mock.Method()
	})

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Method(func() {

		})

		mock.Method()

		fmt.Println(spy)
	})
}
`,
		},

		{
			name: "no arguments, with returns",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				Returns:       varInfos("first: first []string", "second: second error"),
				SkipExpect:    false,
			},
			expected: `package emitter

import (
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock := testTarget()

	t.Run("expect called once", func(t *testing.T) {
		mock.EXPECT().Method(t)

		mock.Method()
	})

	t.Run("expect called twice", func(t *testing.T) {
		mock.EXPECT().Method(t)
		mock.EXPECT().Method(t)

		mock.Method()
		mock.Method()
	})

	t.Run("expect called - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock.EXPECT().Method(t).Return(first, second)

		ret0, ret1 := mock.Method()
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Method(func() ([]string, error) {
			return nil, nil
		})

		mock.Method()

		fmt.Println(spy)
	})
}
`,
		},

		{
			name: "with arguments, no returns",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				Arguments:     varInfos("ctx: ctx context.Context", "tenantID: tenantID string"),
				SkipExpect:    false,
			},
			expected: `package emitter

import (
	"context"
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock := testTarget()

	t.Run("expect called once", func(t *testing.T) {
		mock.EXPECT().Method(t)

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called twice", func(t *testing.T) {
		mock.EXPECT().Method(t)
		mock.EXPECT().Method(t)

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match all arguments by values", func(t *testing.T) {
		var ctx context.Context
		var tenantID string

		mock.EXPECT().Method(t).With(ctx, tenantID)

		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match partial argument by value", func(t *testing.T) {
		var ctx context.Context

		mock.EXPECT().Method(t).WithCtx(ctx)

		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match all arguments by callback", func(t *testing.T) {
		mock.EXPECT().Method(t).Match(func(ctx context.Context, tenantID string) bool {
			return true
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match partial argument by callback", func(t *testing.T) {
		mock.EXPECT().Method(t).MatchCtx(func(ctx context.Context) bool {
			return true
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Method(func(ctx context.Context, tenantID string) {

		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)

		fmt.Println(spy)
	})
}
`,
		},

		{
			name: "with arguments and returns",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				Arguments:     varInfos("ctx: ctx context.Context", "tenantID: tenantID string"),
				Returns:       varInfos("first: first []string", "second: second error"),
				SkipExpect:    false,
			},
			expected: `package emitter

import (
	"context"
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock := testTarget()

	t.Run("expect called once", func(t *testing.T) {
		mock.EXPECT().Method(t)

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called twice", func(t *testing.T) {
		mock.EXPECT().Method(t)
		mock.EXPECT().Method(t)

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock.EXPECT().Method(t).Return(first, second)

		var ctx context.Context
		var tenantID string
		ret0, ret1 := mock.Method(ctx, tenantID)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match all arguments by values", func(t *testing.T) {
		var ctx context.Context
		var tenantID string

		mock.EXPECT().Method(t).With(ctx, tenantID)

		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match all arguments by values - stub return", func(t *testing.T) {
		var ctx context.Context
		var tenantID string
		var first []string
		var second error

		mock.EXPECT().Method(t).With(ctx, tenantID).Return(first, second)

		ret0, ret1 := mock.Method(ctx, tenantID)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match partial argument by value", func(t *testing.T) {
		var ctx context.Context

		mock.EXPECT().Method(t).WithCtx(ctx)

		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match partial by value - stub return", func(t *testing.T) {
		var ctx context.Context
		var first []string
		var second error

		mock.EXPECT().Method(t).WithCtx(ctx).Return(first, second)

		var tenantID string
		ret0, ret1 := mock.Method(ctx, tenantID)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match all arguments by callback", func(t *testing.T) {
		mock.EXPECT().Method(t).Match(func(ctx context.Context, tenantID string) bool {
			return true
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match all arguments by callback - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock.EXPECT().Method(t).Match(func(ctx context.Context, tenantID string) bool {
			return true
		}).Return(first, second)

		var ctx context.Context
		var tenantID string
		ret0, ret1 := mock.Method(ctx, tenantID)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match partial argument by callback", func(t *testing.T) {
		mock.EXPECT().Method(t).MatchCtx(func(ctx context.Context) bool {
			return true
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)
	})

	t.Run("expect called - match partial by callback - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock.EXPECT().Method(t).MatchCtx(func(ctx context.Context) bool {
			return true
		}).Return(first, second)

		var ctx context.Context
		var tenantID string
		ret0, ret1 := mock.Method(ctx, tenantID)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Method(func(ctx context.Context, tenantID string) ([]string, error) {
			return nil, nil
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)

		fmt.Println(spy)
	})
}
`,
		},

		{
			name: "with arguments and returns but skip expect",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				Arguments:     varInfos("ctx: ctx context.Context", "tenantID: tenantID string"),
				Returns:       varInfos("first: first []string", "second: second error"),
				SkipExpect:    true,
			},
			expected: `package emitter

import (
	"context"
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock := testTarget()

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock = testTarget()
		spy := mock.STUB().Method(func(ctx context.Context, tenantID string) ([]string, error) {
			return nil, nil
		})

		var ctx context.Context
		var tenantID string
		mock.Method(ctx, tenantID)

		fmt.Println(spy)
	})
}
`,
		},

		{
			name: "can handle name collision",
			data: ExampleData{
				Constructor:   "testTarget",
				InterfaceName: "Target",
				MethodName:    "Method",
				Arguments:     varInfos("mock: mock context.Context", "spy: spy string"),
				Returns: []VarInfo{
					{Name: "first", Field: "first", OriginalName: "first", Type: gentest.Type("[]string")},
					{Name: "second", Field: "second", OriginalName: "second", Type: gentest.Type("error")},
				},
				SkipExpect: false,
			},
			expected: `package emitter

import (
	"context"
	"fmt"
	"testing"
)

func Test_Target_Method(t *testing.T) {
	mock0 := testTarget()

	t.Run("expect called once", func(t *testing.T) {
		mock0.EXPECT().Method(t)

		var mock context.Context
		var spy string
		mock0.Method(mock, spy)
	})

	t.Run("expect called twice", func(t *testing.T) {
		mock0.EXPECT().Method(t)
		mock0.EXPECT().Method(t)

		var mock context.Context
		var spy string
		mock0.Method(mock, spy)
		mock0.Method(mock, spy)
	})

	t.Run("expect called - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock0.EXPECT().Method(t).Return(first, second)

		var mock context.Context
		var spy string
		ret0, ret1 := mock0.Method(mock, spy)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match all arguments by values", func(t *testing.T) {
		var mock context.Context
		var spy string

		mock0.EXPECT().Method(t).With(mock, spy)

		mock0.Method(mock, spy)
	})

	t.Run("expect called - match all arguments by values - stub return", func(t *testing.T) {
		var mock context.Context
		var spy string
		var first []string
		var second error

		mock0.EXPECT().Method(t).With(mock, spy).Return(first, second)

		ret0, ret1 := mock0.Method(mock, spy)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match partial argument by value", func(t *testing.T) {
		var mock context.Context

		mock0.EXPECT().Method(t).WithMock(mock)

		var spy string
		mock0.Method(mock, spy)
	})

	t.Run("expect called - match partial by value - stub return", func(t *testing.T) {
		var mock context.Context
		var first []string
		var second error

		mock0.EXPECT().Method(t).WithMock(mock).Return(first, second)

		var spy string
		ret0, ret1 := mock0.Method(mock, spy)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match all arguments by callback", func(t *testing.T) {
		mock0.EXPECT().Method(t).Match(func(mock context.Context, spy string) bool {
			return true
		})

		var mock context.Context
		var spy string
		mock0.Method(mock, spy)
	})

	t.Run("expect called - match all arguments by callback - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock0.EXPECT().Method(t).Match(func(mock context.Context, spy string) bool {
			return true
		}).Return(first, second)

		var mock context.Context
		var spy string
		ret0, ret1 := mock0.Method(mock, spy)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("expect called - match partial argument by callback", func(t *testing.T) {
		mock0.EXPECT().Method(t).MatchMock(func(mock context.Context) bool {
			return true
		})

		var mock context.Context
		var spy string
		mock0.Method(mock, spy)
	})

	t.Run("expect called - match partial by callback - stub return", func(t *testing.T) {
		var first []string
		var second error

		mock0.EXPECT().Method(t).MatchMock(func(mock context.Context) bool {
			return true
		}).Return(first, second)

		var mock context.Context
		var spy string
		ret0, ret1 := mock0.Method(mock, spy)
		fmt.Println(ret0, first)
		fmt.Println(ret1, second)
	})

	t.Run("fine-grained control with stub signature", func(t *testing.T) {
		mock0 = testTarget()
		spy0 := mock0.STUB().Method(func(mock context.Context, spy string) (first []string, second error) {
			return nil, nil
		})

		var mock context.Context
		var spy string
		mock0.Method(mock, spy)

		fmt.Println(spy0)
	})
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runEmitterTest(t, &tc.data, tc.expected)
		})
	}
}
