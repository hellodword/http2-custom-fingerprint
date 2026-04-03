// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.27 && http2wrap

// Transport wrapping a net/http.Transport.

package http2

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"
)

func configureTransport(t1 *http.Transport) error {
	// ConfigureTransport is a no-op: The http.Transport already supports HTTP/2.
	return nil
}

func configureTransports(t1 *http.Transport) (*Transport, error) {
	// ConfigureTransport returns an http2.Transport with a configuration
	// linked to the http.Transport's.
	tr2 := &Transport{}
	tr2.configure(t1)
	return tr2, nil
}

// transportConfig is passed to net/http.Transport.RegisterProtocol("http/2", config).
// It provides the net/http.Transport with access to the configuration in the
// x/net/http2.Transport.
type transportConfig struct {
	t *Transport
}

// Registered is called by net/http.Transport.RegisterProtocol,
// to let us know that it understands the registration mechanism we're using.
func (t transportConfig) Registered(t1 *http.Transport) {
	t.t.t1 = t1
}

func (t transportConfig) DisableCompression() bool {
	return t.t.DisableCompression
}

func (t transportConfig) MaxHeaderListSize() int64 {
	return int64(t.t.MaxHeaderListSize)
}

func (t transportConfig) IdleConnTimeout() time.Duration {
	return t.t.IdleConnTimeout
}

func (t transportConfig) HTTP2Config() http.HTTP2Config {
	return http.HTTP2Config{
		StrictMaxConcurrentRequests: t.t.StrictMaxConcurrentStreams,
		MaxDecoderHeaderTableSize:   int(t.t.MaxDecoderHeaderTableSize),
		MaxEncoderHeaderTableSize:   int(t.t.MaxEncoderHeaderTableSize),
		MaxReadFrameSize:            int(t.t.MaxReadFrameSize),
		SendPingTimeout:             t.t.ReadIdleTimeout,
		PingTimeout:                 t.t.PingTimeout,
		WriteByteTimeout:            t.t.WriteByteTimeout,
		CountError:                  t.t.CountError,
	}
}

// ExternalRoundTrip reports whether the Transport wants to take control of the RoundTrip call.
// If the user hasn't configured a custom connection pool, we leave the RoundTrip up to net/http.
func (t transportConfig) ExternalRoundTrip() bool {
	return t.t.ConnPool != nil
}

// RoundTrip is used when the http.Transport is passing control of the full
// RoundTrip to us--connection pooling, retries, etc.
//
// This is only used when the http2.Transport has a user-provided ConnPool.
// Any other time, net/http handles everything.
func (t transportConfig) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.t.ConnPool == nil {
		return nil, http.ErrSkipAltProtocol
	}
	return t.t.RoundTrip(req)
}

// netConnContextKey passes a net.Conn to http.Transport.NewClientConn.
// See http2.Transport.NewClientConn.
type netConnContextKey struct{}

// ConnFromContext lets the http.Transport fetch a net.Conn out of a context
// passed to NewClientConn. See http2.Transport.NewClientConn.
func (t transportConfig) ConnFromContext(ctx context.Context) net.Conn {
	nc, _ := ctx.Value(netConnContextKey{}).(net.Conn)
	return nc
}

// http2TransportContextKey marks a RoundTrip as needing its dial handled by the http2.Transport.
// We set this for http2.RoundTrip calls, where the historical behavior is to use the
// http2.Transport's dialer.
type http2TransportContextKey struct{}

// DialFromContext dials a new connection using the http2.Transport's DialTLS/DialTLSContext.
func (t transportConfig) DialFromContext(ctx context.Context, network, address string) (net.Conn, error) {
	if ctx.Value(http2TransportContextKey{}) == nil {
		// We're being called from a RoundTrip that did not start with an http2.Transport.
		// Use the http.Transport's dialer.
		return nil, errors.ErrUnsupported
	}

	tlsConf := t.t.TLSClientConfig
	if tlsConf == nil {
		tlsConf = &tls.Config{}
	} else {
		tlsConf = tlsConf.Clone()
	}
	if !slices.Contains(tlsConf.NextProtos, "h2") {
		tlsConf.NextProtos = append([]string{"h2"}, tlsConf.NextProtos...)
	}
	if tlsConf.ServerName == "" {
		host, _, err := net.SplitHostPort(address)
		if err == nil {
			tlsConf.ServerName = host
		}
	}
	return t.t.dialTLS(ctx, network, address, tlsConf)
}

type transportInternal struct {
	initOnce sync.Once
	t1       *http.Transport
}

func (t *Transport) init() {
	t.initOnce.Do(func() {
		if t.t1 != nil {
			return
		}
		t1 := http.DefaultTransport.(*http.Transport).Clone()
		t.configure(t1)
	})
}

func (t *Transport) configure(t1 *http.Transport) {
	t1.RegisterProtocol("http/2", transportConfig{t})
	// tr2.t1 is set by transportConfig.Registered.
	if t.t1 != t1 {
		panic("http2: net/http does not support this version of x/net/http2")
	}
}

func (t *Transport) roundTripOpt(req *http.Request, opt RoundTripOpt) (*http.Response, error) {
	t.init()

	if req.URL.Scheme == "http" && !t.AllowHTTP {
		return nil, errors.New("http2: unencrypted HTTP/2 not enabled")
	}

	// When the Transport has a user-provided connection pool (unusual, deprecated),
	// we need to handle picking a connection, retrys, etc.
	if t.ConnPool != nil {
		return nil, errors.New("TODO")
	}

	// Setting this context key lets net/http know that if it is necessary to dial
	// a new connection, we should handle the net.Dial.
	//
	// Both http.Transport and http2.Transport allow the user to provide a custom
	// dial function, and historically you only get the dial function from the
	// Transport you're calling RoundTrip on.
	ctx := context.WithValue(req.Context(), http2TransportContextKey{}, t)
	req = req.WithContext(ctx)

	return t.t1.RoundTrip(req)
}

func (t *Transport) closeIdleConnections() {
	t.init()
	t.t1.CloseIdleConnections()
}

func (t *Transport) newUserClientConn(c net.Conn) (*ClientConn, error) {
	return nil, errors.New("TODO")
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
