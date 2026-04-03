// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.27 && http2wrap

// Wrapping the HTTP/2 implementation in net/http/internal/http2:
// Internal package hooks not available.

package http2_test

import (
	"errors"
	"net/http"
	"testing"
)

const wrappedAPI = true

type httpClientConn = http.ClientConn

func newTestClientConn(t testing.TB, opts ...any) *testClientConn {
	t.Fatal("TODO")
	return nil
}

func (tr *testTransport) maybeAddNewClientConnHook() {
}

func (tc *testClientConn) doRoundTrip(req *http.Request, f func(streamID uint32)) (*http.Response, error) {
	return nil, errors.New("TODO")
}
