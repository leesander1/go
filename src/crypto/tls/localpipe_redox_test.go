// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package tls

import (
	"net"
	"os"
	"syscall"
	"testing"
)

func localPipeOS(t testing.TB) (net.Conn, net.Conn, bool) {
	// Redox's loopback listener can accept too slowly or report unstable
	// socket addresses in this helper. Use a buffered socketpair instead.
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("localPipe: socketpair: %v", err)
	}

	f1 := os.NewFile(uintptr(fds[0]), "localPipe-1")
	f2 := os.NewFile(uintptr(fds[1]), "localPipe-2")
	defer f1.Close()
	defer f2.Close()

	c1, err := net.FileConn(f1)
	if err != nil {
		t.Fatalf("localPipe: file conn 1: %v", err)
	}
	c2, err := net.FileConn(f2)
	if err != nil {
		c1.Close()
		t.Fatalf("localPipe: file conn 2: %v", err)
	}

	t.Cleanup(func() { c1.Close() })
	t.Cleanup(func() { c2.Close() })
	return c1, c2, true
}
