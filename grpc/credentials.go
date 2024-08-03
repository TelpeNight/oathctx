package grpc

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/TelpeNight/oauthctx"
)

// TokenSource supplies PerRPCCredentials from an oauth2.TokenSource.
type TokenSource struct {
	oauthctx.TokenSource
}

var _ credentials.PerRPCCredentials = (*TokenSource)(nil)

// GetRequestMetadata gets the request metadata as a map from a TokenSource.
func (ts *TokenSource) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := ts.TokenContext(ctx)
	if err != nil {
		return nil, err
	}
	// reusing grpc module impl
	staticSource := oauth2.StaticTokenSource(token)
	impl := oauth.TokenSource{TokenSource: staticSource}
	return impl.GetRequestMetadata(ctx, uri...)
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (ts *TokenSource) RequireTransportSecurity() bool {
	return oauth.TokenSource{}.RequireTransportSecurity()
}
