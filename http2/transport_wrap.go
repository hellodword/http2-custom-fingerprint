// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.27 && http2wrap

// Transport wrapping a net/http.Transport.

package http2

import (
	"context"
	"net"
	"net/http"
)

func configureTransport(t1 *http.Transport) error {
	return nil
}

func configureTransports(t1 *http.Transport) (*Transport, error) {
	tr2 := &Transport{}
	return tr2, nil
}

type transportInternal struct {
}

func (t *Transport) roundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	return nil, nil
}

func (t *Transport) closeIdleConnections() {
}

func (t *Transport) newUserClientConn(c net.Conn) (*ClientConn, error) {
	return nil, nil
}

// ClientConn is the state of a single HTTP/2 client connection to an
// HTTP/2 server.
type ClientConn struct {
	atomicReused uint32 // whether conn is being reused; atomic
}

func (cc *ClientConn) roundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (cc *ClientConn) canTakeNewRequest() bool {
	return false
}

func (cc *ClientConn) close() error {
	return nil
}

func (cc *ClientConn) ping(ctx context.Context) error {
	return nil
}

func (cc *ClientConn) reserveNewRequest() bool {
	return false
}

func (cc *ClientConn) setDoNotReuse() {
}

func (cc *ClientConn) shutdown(ctx context.Context) error {
	return nil
}

func (cc *ClientConn) state() ClientConnState {
	return ClientConnState{}
}

func (cc *ClientConn) stopIdleTimer() {}

func traceGotConn(req *http.Request, cc *ClientConn, reused bool) {
}
