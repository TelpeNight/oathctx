package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type ConfigClientOp func(o *configClientOp)

func ConfigClientWithTransportClient(client *http.Client) ConfigClientOp {
	return func(o *configClientOp) {
		o.authClient = client
		o.requestClient = client
	}
}

func ConfigClientWithAuthClient(client *http.Client) ConfigClientOp {
	return func(o *configClientOp) {
		o.authClient = client
	}
}

func ConfigClientWithRequestClient(client *http.Client) ConfigClientOp {
	return func(o *configClientOp) {
		o.requestClient = client
	}
}

type TokenSourceOp func(o *tokenSourceOps)

func TokenSourceWithClient(client *http.Client) TokenSourceOp {
	return func(o *tokenSourceOps) {
		o.client = client
	}
}

// Oauth2ContextClient may be used with options. Returns nil on missing value
func Oauth2ContextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	return nil
}

type configClientOp struct {
	requestClient *http.Client
	authClient    *http.Client
}

func (o *configClientOp) tokenSourceOps() []TokenSourceOp {
	if o.authClient == nil {
		return nil
	}
	return []TokenSourceOp{TokenSourceWithClient(o.authClient)}
}

func (o *configClientOp) clientOps() []ClientOp {
	if o.requestClient == nil {
		return nil
	}
	return []ClientOp{ClientWithRequestClient(o.requestClient)}
}

type tokenSourceOps struct {
	client *http.Client
}

func makeTokenSourceOps(ops []TokenSourceOp) *tokenSourceOps {
	var options *tokenSourceOps
	if len(ops) > 0 {
		options = &tokenSourceOps{}
		for _, op := range ops {
			op(options)
		}
	}
	return options
}

func (o *tokenSourceOps) ctx(ctx context.Context) context.Context {
	if o == nil || o.client == nil {
		return ctx
	}
	return context.WithValue(ctx, oauth2.HTTPClient, o.client)
}
