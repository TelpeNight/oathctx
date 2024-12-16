# oauthctx

---

This module is used to bypass golang oauth2 package token source limitations: https://github.com/golang/oauth2/issues/262

At current moment, we can't pass a request context to an oauth request, so it doesn't respect deadlines and values.

This package reimplements only small subset of functionality and should play well with any library.

The goal of this package is not to reimplement oauth2 details, but focus on interfaces. So we are trying to reuse existing
implementation as much as possible.

The algorithm is the following:

1. Wrap existing functionality into `TokenSource` implementation.
This may come in to different flavours:

    * 'Immutable' token sources. They don't have any mutable inner state. These sources can be converted as simple as following:
        ```go
        src := c.new.TokenSource(ctx)
        return src.Token()
        ```
      See [convert.go](convert.go) and [clientcredentials.go](clientcredentials.go) for the reference.
      Package also provides generic purpose `ConvertImmutable` and `NewOauth2TokenSource` to simplify this.
   
    * 'Mutable' token sources, such as `oauth2.tokenRefresher`. They have inner state, that can be updated on any call.
      To mimic this behavior we need some extra work. See [config.go](config.go)

2. Then wrap new `oauthctx.TokenSource` with `oauthctx.ReuseTokenSource`. It will refresh expired token.
   Also it provides Context-aware synchronisation.

3. oauth2.HTTPClient and similar functionality is provided by Options. Using Context to achieve this is sooo messy.
   To be closer to original functionality, http.client which was provided on construction has higher priority over http.client
   from request's context (if there is one). But in general, you can provide per-call ctx with custom client and other values. 
   This behavior is different from the original library.

And that is it. Any existing `TokenSource` that holds internal `Context` may be used in this way.

The next "a-ha" moment is that we can reuse any other module, which depends on `ouath2`, using `oauth2.StaticTokenSource`.
Just obtain token with `TokenContext` method. Then use existing implementation with `StaticTokenSource`.
No need to dive in implementation details.
See [transport.go](transport.go) and [grpc/credentials.go](grpc/credentials.go) for the reference.

Currently, module provides context-aware implementation of `http.RoundTripper` and `grpc.PerRPCCredentials`. Feel free to pull request new ones.

## Code examples

---

```go
// grpc
package main
import (
    "golang.org/x/oauth2"
	
    "google.golang.org/grpc/credentials"
    gcred "google.golang.org/grpc/credentials/google"

    "github.com/TelpeNight/oauthctx"
    grpcctx "github.com/TelpeNight/oauthctx/grpc"
)

var conf = oauthctx.NewConfig(&oauth2.Config{
    //...
})
var refreshToken string = "..."

ts := conf.TokenSource(
	&oauth2.Token{RefreshToken: refreshToken},
	// custom http.Client can be provided with option
	oauthctx.TokenSourceWithClient(...))
ts = oauthctx.ReuseTokenSource(nil, ts)

var bundle credentials.Bundle = gcred.NewDefaultCredentialsWithOptions(
    gcred.DefaultCredentialsOptions{
        PerRPCCreds: &grpcctx.TokenSource{
            TokenSource: ts,
        },
    },
)

// use bundle to create a client. methods' context will be passed to oauth2, so overall call will respect timeouts
```
