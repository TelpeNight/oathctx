# oathctx

---

This module is used to bypass golang oauth2 package token source limitations: https://github.com/golang/oauth2/issues/262

At current moment, we can't pass a request context to an oauth request, so it doesn't respect deadlines and values.

The main idea is using ctxref/ContextReference, which lets caller of TokenSource to switch internal ctx before calling Token() method.

This package reimplements only small subset of functionality and should play well with any library.

The main strategy is to convert oauth2.TokenSource to oauthctx.TokenSource via ReuseTokenSource method. Then we should switch to context-aware implementation of TokenSource caller.

And that is it. Any existing TokenSource that holds internal Context may be used in this way. (See [Thread safety](#thread-safety-notes) for current limitations)

The next "a-ha" moment is that we can reuse any existing ouath2 client with oauth2.StaticTokenSource. Just obtain token with TokenContext method. Then use basic implementation with StaticTokenSource. No need to dive in implementation details. 

Currently, module provides context-aware implementation of http.RoundTripper and grpc.PerRPCCredentials. Feel free to pull request new ones.

## Code examples

---

```go
// OAuth transport
package main
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/clientcredentials"

    "github.com/TelpeNight/oauthctx"
    "github.com/TelpeNight/oauthctx/ctxref"
)


ctx := ctxref.Background()

// can be used as ordinary context:
transportCtx := context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
    Transport: unwrapTokenTransport,
})

var config clientcredentials.Config = ...
transport = &oauthctx.Transport{
    // note ctx is passed to ReuseTokenSource, and its child to TokenSource
    Source: oauthctx.ReuseTokenSource(ctx, nil, config.TokenSource(transportCtx)),
}

client := &http.Client{
    Transport:  transport,
}

---

reqCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
req := http.NewRequestWithContext(reqCtx, ...)
client.Do(req) // <- oauth call will respect timeout
```

```go
// GRPC
package main
import (
    "golang.org/x/oauth2"
	
    "google.golang.org/grpc/credentials"
    gcred "google.golang.org/grpc/credentials/google"

    "github.com/TelpeNight/oauthctx"
    "github.com/TelpeNight/oauthctx/ctxref"
    grpcctx "github.com/TelpeNight/oauthctx/grpc"
)

ctx := ctxref.Background()

var conf = &oauth2.Config{
    ...
}
var refreshToken string = ...

classicSrc := conf.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
ts := oauthctx.ReuseTokenSource(ctx, nil, classicSrc)

var bundle credentials.Bundle = gcred.NewDefaultCredentialsWithOptions(
    gcred.DefaultCredentialsOptions{
        PerRPCCreds: &grpcctx.TokenSource{
            TokenSource: ts,
        },
    },
)

// use bundle to create a client. methods' context will be passed to oauth2, so overall call will respect timeouts

```


## Thread safety notes

---

By default `ctxref.ContextReference` is not thread safe and is guarded by `ReuseTokenSource`.
If you need thread safe version - consider making a pull request. But keep in mind, that synchronisation should respect ctx's timeouts.
See `ReuseTokenSource` for a reference. Trivial Mutex implementations won't be accepted.

After calling `ctx.Use(other)` ctx becomes a "perfect forwarder" to other context. 
At the moment of calling to `ctx.Unuse()` ctx should not be accessed by any concurrent goroutine.
So, underling implementation should not retain context by any means other than storing it inside `oauth2.TokenSource`.
Concurrent access to ctx while calling to `Use` or `Unuse` may cause serious problem!  

It seems that all supported modules are safe-for-use in this manner.
All them use `golang.org/x/oauth2/internal.RetrieveToken`.
Which does `ContextClient(ctx).Do(req.WithContext(ctx))`.
I don't expect any issues with std library after all calls have returned.
You may want to check custom `ContextClient` if you, your platform or tools use one.

But if you'll encounter ctxref leaking - please fire an issue.
