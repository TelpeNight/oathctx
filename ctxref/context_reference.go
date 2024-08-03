package ctxref

import "context"

// ContextReference default implementations are not thread-safe and usually guarded by oauthctx.ReuseTokenSource
type ContextReference interface {
	context.Context
	Use(ctx context.Context)
	Unuse()
}

func Background() ContextReference {
	return &backgroundContextReference{context.Background()}
}

type backgroundContextReference struct {
	context.Context
}

func (c *backgroundContextReference) Use(ctx context.Context) {
	c.Context = ctx
}

func (c *backgroundContextReference) Unuse() {
	c.Context = context.Background()
}
