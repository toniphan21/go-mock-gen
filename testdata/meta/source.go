package meta

import "context"

type Result struct{}

type Target interface {
	//Empty()
	//NoReturn(ctx context.Context, input string)
	//NoArgument() ([]Result, error)

	Full(ctx context.Context, input string) ([]Result, error)
}
