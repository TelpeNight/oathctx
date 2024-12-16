package oauthctx

import (
	"context"

	"golang.org/x/oauth2"
)

type TokenSource interface {
	TokenContext(ctx context.Context) (*oauth2.Token, error)
}

// ReuseTokenSource returns a TokenSource which repeatedly returns the
// same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src.
func ReuseTokenSource(t *oauth2.Token, src TokenSource) TokenSource {
	if t == nil {
		if rt, ok := src.(*reuseTokenSource); ok {
			//if t == nil {
			// Just use it directly.
			return rt
			//}
			// not secure, as rt.new is not guarded by single mutex anymore
			//src = rt.new
		}
	}
	ts := &reuseTokenSource{
		new: src,
		mu:  make(chan struct{}, 1),
		t:   t,
	}
	ts.mu <- struct{}{}
	return ts
}

type reuseTokenSource struct {
	new TokenSource   // called when t is expired.
	mu  chan struct{} // guards t. use chain to select on mu and context
	t   *oauth2.Token
}

func (s *reuseTokenSource) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	select {
	case <-s.mu:
		defer func() { s.mu <- struct{}{} }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if s.t.Valid() {
		return s.t, nil
	}

	t, err := s.new.TokenContext(ctx)
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}
