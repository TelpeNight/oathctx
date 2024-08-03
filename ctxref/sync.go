package ctxref

import (
	"context"
	"sync"
)

func Sync(ctx ContextReference) ContextReference {
	if _, is := ctx.(*syncContextRef); is {
		return ctx
	}
	return &syncContextRef{ContextReference: ctx}
}

type syncContextRef struct {
	ContextReference
	mu sync.Mutex
}

func (s *syncContextRef) Use(ctx context.Context) {
	s.mu.Lock()
	s.ContextReference.Use(ctx)
}

func (s *syncContextRef) Unuse() {
	s.ContextReference.Unuse()
	s.mu.Unlock()
}
