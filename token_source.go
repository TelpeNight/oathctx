package oauthctx

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/TelpeNight/oauthctx/ctxref"
)

type TokenSource interface {
	TokenContext(ctx context.Context) (*oauth2.Token, error)
}

// ReuseTokenSource returns a TokenSource which repeatedly returns the
// same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src.
// oauth2.TokenSource should be seeded with ctx
func ReuseTokenSource(ctx ctxref.ContextReference, t *oauth2.Token, src oauth2.TokenSource) TokenSource {
	ts := &reuseTokenSource{
		ctx: ctx,
		t:   t,
		new: src,
		mu:  make(chan struct{}, 1),
	}
	ts.mu <- struct{}{}
	return ts
}

type reuseTokenSource struct {
	ctx ctxref.ContextReference // ctx of new
	new oauth2.TokenSource      // called when t is expired.

	mu chan struct{} // guards t and ctx. use chain to select on mu and context
	t  *oauth2.Token
}

func (s *reuseTokenSource) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	select {
	case l := <-s.mu:
		defer func() {
			s.mu <- l
		}()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if s.t.Valid() {
		return s.t, nil
	}

	s.ctx.Use(ctx)
	defer s.ctx.Unuse()

	// s.new is using s.ctx if everything is set up correctly
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}
