package oauthctx

import (
	"context"

	"golang.org/x/oauth2"
)

type Oauth2TokenConfig interface {
	TokenSource(ctx context.Context) oauth2.TokenSource
}

func ConvertImmutable(conf Oauth2TokenConfig, ops ...TokenSourceOp) TokenSource {
	return &convertTokenSrc{
		new: conf,
		ops: makeTokenSourceOps(ops),
	}
}

type convertTokenSrc struct {
	new Oauth2TokenConfig
	ops *tokenSourceOps
}

func (c *convertTokenSrc) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	src := c.new.TokenSource(c.ops.ctx(ctx))
	return src.Token()
}

type NewOauth2TokenSource func(ctx context.Context) oauth2.TokenSource

func (n NewOauth2TokenSource) TokenSource(ctx context.Context) oauth2.TokenSource {
	return n(ctx)
}
