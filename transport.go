package oauthctx

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// Transport is a http.RoundTripper that makes OAuth 2.0 HTTP requests,
// wrapping a base RoundTripper and adding an Authorization header
// with a token from the supplied Sources.
//
// Transport is a low-level mechanism. Most code will use the
// higher-level NewClient function instead.
type Transport struct {
	// Source supplies the token to add to outgoing requests'
	// Authorization headers.
	Source TokenSource

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with an
// access token from Transport's Source.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Source == nil {
		return nil, errors.New("oauth2: Transport's Source is nil")
	}
	token, err := t.Source.TokenContext(req.Context())
	if err != nil {
		return nil, err
	}

	// reusing oauth2 module impl
	staticSource := oauth2.StaticTokenSource(token)
	impl := oauth2.Transport{
		Source: staticSource,
		Base:   t.Base,
	}
	return impl.RoundTrip(req)
}
