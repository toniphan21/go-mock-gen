package meta

import "context"

type Production struct {
	target Target
}

func (p *Production) CallFullOnce(ctx context.Context, input string) {
	// do something awesome

	// do something else more awesome

	_, _ = p.target.Full(ctx, input)
}

func (p *Production) CallFullTwice(ctx context.Context, input string) {
	_, _ = p.target.Full(ctx, input+" 1")

	// do something awesome

	// do something else more awesome

	_, _ = p.target.Full(ctx, input+" 2")
}

func (p *Production) CallFullThrice(ctx context.Context, input string) {
	_, _ = p.target.Full(ctx, input+" 1")

	// do something awesome

	// do something else more awesome

	_, _ = p.target.Full(ctx, input+" 2")

	// legend...

	// wait for it

	// ...ary
	_, _ = p.target.Full(ctx, input+" 3")

	// never stop partying - uncle Jerry said
}
