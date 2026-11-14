// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"syscall"
	"testing"
)

func socketBufferSize(t *testing.T, c syscall.Conn, option int) int {
	t.Helper()
	raw, err := c.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	var value int
	var optionErr error
	if err := raw.Control(func(fd uintptr) {
		value, optionErr = syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, option)
	}); err != nil {
		t.Fatal(err)
	}
	if optionErr != nil {
		t.Fatal(optionErr)
	}
	return value
}

func testSocketBufferOptions(t *testing.T, name string, c interface {
	syscall.Conn
	SetReadBuffer(int) error
	SetWriteBuffer(int) error
}) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		const requested = 2048
		if err := c.SetReadBuffer(requested); err != nil {
			t.Fatalf("SetReadBuffer(%d): %v", requested, err)
		}
		if got := socketBufferSize(t, c, syscall.SO_RCVBUF); got < requested {
			t.Errorf("SO_RCVBUF = %d, want at least %d", got, requested)
		}
		if err := c.SetWriteBuffer(requested); err != nil {
			t.Fatalf("SetWriteBuffer(%d): %v", requested, err)
		}
		if got := socketBufferSize(t, c, syscall.SO_SNDBUF); got < requested {
			t.Errorf("SO_SNDBUF = %d, want at least %d", got, requested)
		}
		if err := c.SetReadBuffer(0); err == nil {
			t.Error("SetReadBuffer(0) succeeded, want error")
		}
		if err := c.SetWriteBuffer(-1); err == nil {
			t.Error("SetWriteBuffer(-1) succeeded, want error")
		}
	})
}

func TestRedoxSocketBufferOptions(t *testing.T) {
	ln, err := ListenTCP("tcp4", &TCPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

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
	testSocketBufferOptions(t, "tcp", accepted)

	select {
	case err := <-dialErr:
		t.Fatal(err)
	case conn := <-dialed:
		defer conn.Close()
	}

	udp, err := ListenUDP("udp4", &UDPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer udp.Close()
	testSocketBufferOptions(t, "udp", udp)
}
