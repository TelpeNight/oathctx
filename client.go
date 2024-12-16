package oauthctx

import (
	"net/http"
)

// NewClient creates an *http.Client from TokenSource.
func NewClient(src TokenSource, ops ...ClientOp) *http.Client {
	var options clientOp
	for _, op := range ops {
		op(&options)
	}
	if options.client == nil {
		options.client = http.DefaultClient
	}
	if src == nil {
		// to match original behaviour
		return options.client
	}
	return &http.Client{
		Transport: &Transport{
			Source: ReuseTokenSource(nil, src),
			Base:   options.client.Transport,
		},
	}
}

// ClientWithRequestClient sets transport, that will be used after obtaining token to make a real request
func ClientWithRequestClient(client *http.Client) ClientOp {
	return func(o *clientOp) {
		o.client = client
	}
}

type ClientOp func(o *clientOp)

type clientOp struct {
	client *http.Client
}
