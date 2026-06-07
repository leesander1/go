// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package tls

import (
	"net"
	"testing"
)

func localPipeOS(t testing.TB) (net.Conn, net.Conn, bool) {
	return nil, nil, false
}
