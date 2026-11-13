// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"net/netip"
	"testing"
)

func TestDialTCPAddrFromAddrPort(t *testing.T) {
	if addr := dialTCPAddrFromAddrPort(netip.AddrPort{}); addr != nil {
		t.Fatalf("zero address converted to %#v; want nil", addr)
	}

	want := netip.MustParseAddrPort("127.0.0.1:1234")
	if addr := dialTCPAddrFromAddrPort(want); addr == nil || addr.AddrPort() != want {
		t.Fatalf("valid address converted to %#v; want %v", addr, want)
	}
}
