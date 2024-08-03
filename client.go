package oauthctx

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/TelpeNight/oauthctx/ctxref"
)

// NewClient creates an *http.Client from a Context and oauth2.TokenSource.
// The returned client is not valid beyond the lifetime of the context.
// oauth2.TokenSource should be seeded with ctx.
//
//	ctx := ctxref.Background()
//	transportCtx := context.WithValue(ctx, oauth2.HTTPClient, customTokenClient)
//	oauthClient := oauthctx.NewClient(ctx, tokenSource.TokenSource(transportCtx))
func NewClient(ctx ctxref.ContextReference, src oauth2.TokenSource) *http.Client {
	return &http.Client{
		Transport: &Transport{
			Source: ReuseTokenSource(ctx, nil, src),
		},
	}
}
