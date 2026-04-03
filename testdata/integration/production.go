package integration

import "context"

type createUserOp struct {
	repo Repository
}

func (op *createUserOp) Handle(ctx context.Context, user User) error {
	return op.repo.CreateUser(ctx, user)
}
