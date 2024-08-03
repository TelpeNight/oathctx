# oathctx

This module is used to bypass golang oauth2 package token source limitations: https://github.com/golang/oauth2/issues/262

In current state, we can't pass a request context to an oauth request, so it doesn't respect deadlines and values.

The main idea is using ctxref/ContextReference, which lets caller of TokenSource to switch internal ctx before calling Token() method.

This package reimplements only small subset of functionality and should play well with any library.

The main strategy is to convert oauth2.TokenSource to oauthctx.TokenSource via ReuseTokenSource method. Then we should switch to context-aware implementation of TokenSource caller.

And that is it. Any existing TokenSource that holds internal Context may be used in this way.

The next "a-ha" moment is that we can reuse any trivial caller implementation with oauth2.StaticTokenSource. Just obtain token with TokenContext method. Then use basic implementation with StaticTokenSource. No need to dive in implementation details. 

Currently, module provides context-aware implementation of http.RoundTripper and grpc.PerRPCCredentials. Feel free to pull other cases. 