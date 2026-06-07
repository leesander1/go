// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import "io"

type redoxTCPConnWithoutReadFrom struct {
	c *TCPConn
}

func (c redoxTCPConnWithoutReadFrom) Write(p []byte) (int, error) {
	return c.c.Write(p)
}

func tcpReadFromConn(c *TCPConn) io.Writer {
	return redoxTCPConnWithoutReadFrom{c: c}
}

func tcpReadFromReader(r io.Reader) io.Reader {
	return struct{ io.Reader }{r}
}
