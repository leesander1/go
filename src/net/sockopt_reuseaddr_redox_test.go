// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"errors"
	"syscall"
	"testing"
)

func setSocketOptionInt(t *testing.T, c syscall.Conn, option, value int) error {
	t.Helper()
	raw, err := c.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	var optionErr error
	if err := raw.Control(func(fd uintptr) {
		optionErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, option, value)
	}); err != nil {
		t.Fatal(err)
	}
	return optionErr
}

func TestRedoxTCPReuseAddr(t *testing.T) {
	ln, err := ListenTCP("tcp4", &TCPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	if got := socketOptionInt(t, ln, syscall.SO_REUSEADDR); got != 1 {
		t.Fatalf("listener SO_REUSEADDR = %d, want 1", got)
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

	if got := socketOptionInt(t, accepted, syscall.SO_REUSEADDR); got != 1 {
		t.Fatalf("accepted SO_REUSEADDR = %d, want 1", got)
	}
	if err := setSocketOptionInt(t, accepted, syscall.SO_REUSEADDR, 0); err != nil {
		t.Fatalf("disable SO_REUSEADDR: %v", err)
	}
	if got := socketOptionInt(t, accepted, syscall.SO_REUSEADDR); got != 0 {
		t.Fatalf("disabled SO_REUSEADDR = %d, want 0", got)
	}

	const unsupportedOption = 0x7fffffff
	if err := setSocketOptionInt(t, accepted, unsupportedOption, 1); !errors.Is(err, syscall.ENOPROTOOPT) {
		t.Fatalf("unsupported socket option error = %v, want %v", err, syscall.ENOPROTOOPT)
	}
}
