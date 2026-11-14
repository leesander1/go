// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"syscall"
	"testing"
)

func socketAcceptConn(t *testing.T, c syscall.Conn) int {
	t.Helper()
	raw, err := c.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	var value int
	var optionErr error
	if err := raw.Control(func(fd uintptr) {
		value, optionErr = syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_ACCEPTCONN)
	}); err != nil {
		t.Fatal(err)
	}
	if optionErr != nil {
		t.Fatal(optionErr)
	}
	return value
}

func TestSocketAcceptConn(t *testing.T) {
	ln, err := ListenTCP("tcp4", &TCPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	if got := socketAcceptConn(t, ln); got != 1 {
		t.Fatalf("listener SO_ACCEPTCONN = %d, want 1", got)
	}

	dialed := make(chan *TCPConn, 1)
	dialErr := make(chan error, 1)
	go func() {
		conn, err := DialTCP("tcp4", nil, ln.Addr().(*TCPAddr))
		if err != nil {
			dialErr <- err
			return
		}
		dialed <- conn
	}()

	accepted, err := ln.AcceptTCP()
	if err != nil {
		t.Fatal(err)
	}
	defer accepted.Close()

	select {
	case err := <-dialErr:
		t.Fatal(err)
	case conn := <-dialed:
		defer conn.Close()
	}

	if got := socketAcceptConn(t, accepted); got != 0 {
		t.Fatalf("accepted connection SO_ACCEPTCONN = %d, want 0", got)
	}
}
