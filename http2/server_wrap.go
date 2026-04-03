// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.27 && http2wrap

// Server wrapping a net/http.Server.

package http2

import (
	"errors"
	"net"
	"net/http"
)

type serverInternalState struct {
}

func configureServer(s *http.Server, conf *Server) error {
	return errors.New("TODO")
}

func (s *Server) serveConn(c net.Conn, opts *ServeConnOpts, newf func(*ServerConn)) {
	c.Close() // TODO
}

// FrameWriteRequest is a request to write a frame.
//
// Deprecated: User-provided write schedulers are deprecated.
type FrameWriteRequest struct {
	// Ideally we'd define this in writesched_common.go,
	// to avoid duplicating an exported symbol across two files,
	// but the changes required to make this work are fairly large.
}
