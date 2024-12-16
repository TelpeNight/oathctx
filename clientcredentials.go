package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientCredentials describes a 2-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type ClientCredentials struct{ *clientcredentials.Config }

func NewClientCredentials(cfg *clientcredentials.Config) *ClientCredentials {
	return &ClientCredentials{cfg}
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary.
//
// The provided options control which HTTP client
// is used.
//
// The returned Client and its Transport should not be modified.
func (c *ClientCredentials) Client(ops ...ConfigClientOp) *http.Client {
	var options configClientOp
	for _, op := range ops {
		op(&options)
	}
	return NewClient(
		c.tokenSource(options.tokenSourceOps()), // NewClient will reuse tokenSource
		options.clientOps()...)
}

// Token uses client credentials to retrieve a token.
//
// The provided options optionally controls which HTTP client is used.
func (c *ClientCredentials) Token(ctx context.Context, ops ...TokenSourceOp) (*oauth2.Token, error) {
	return c.tokenSource(ops).TokenContext(ctx)
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided options and the
// client ID and client secret.
//
// Most users will use Config.Client instead.
func (c *ClientCredentials) TokenSource(ops ...TokenSourceOp) TokenSource {
	return ReuseTokenSource(nil, c.tokenSource(ops))
}

func (c *ClientCredentials) tokenSource(ops []TokenSourceOp) TokenSource {
	return ConvertImmutable(c.Config, ops...)
}
