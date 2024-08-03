package ctxref

import "context"

func Identity(ctx context.Context) ContextReference {
	return &identity{ctx}
}

type identity struct {
	context.Context
}

func (i *identity) Use(context.Context) {}

func (i *identity) Unuse() {}
