// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"context"
	"testing"
)

func TestDialerExplicitLocalAddr(t *testing.T) {
	ln, err := ListenTCP("tcp4", &TCPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	accepted := make(chan error, 1)
	go func() {
		conn, err := ln.AcceptTCP()
		if err == nil {
			err = conn.Close()
		}
		accepted <- err
	}()

	wantIP := IPv4(127, 0, 0, 1)
	dialer := Dialer{LocalAddr: &TCPAddr{IP: wantIP}}
	conn, err := dialer.DialContext(context.Background(), "tcp4", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	local := conn.LocalAddr().(*TCPAddr)
	if !local.IP.Equal(wantIP) || local.Port == 0 {
		t.Errorf("local address = %v; want %v with an ephemeral port", local, wantIP)
	}
	if err := conn.Close(); err != nil {
		t.Error(err)
	}
	if err := <-accepted; err != nil {
		t.Error(err)
	}
}
