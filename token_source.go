package oauthctx

import (
	"context"
	"sync"

	"golang.org/x/oauth2"

	"github.com/TelpeNight/oauthctx/ctxref"
)

type TokenSource interface {
	TokenContext(ctx context.Context) (*oauth2.Token, error)
}

// ReuseTokenSource returns a TokenSource which repeatedly returns the
// same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src.
func ReuseTokenSource(ctx ctxref.ContextReference, t *oauth2.Token, src oauth2.TokenSource) TokenSource {
	return &reuseTokenSource{
		ctx: ctx,
		t:   t,
		new: src,
	}
}

type reuseTokenSource struct {
	ctx ctxref.ContextReference // ctx of new
	new oauth2.TokenSource      // called when t is expired.

	mu sync.Mutex // guards t and ctx
	t  *oauth2.Token
}

func (s *reuseTokenSource) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}

	s.ctx.Use(ctx)
	defer s.ctx.Unuse()
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}
