// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall_test

import (
	"syscall"
	"testing"
)

func TestMmapReturnsErrorOnFailure(t *testing.T) {
	b, err := syscall.Mmap(12345, 0, syscall.Getpagesize(), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err == nil {
		_ = syscall.Munmap(b)
		t.Fatalf("Mmap with invalid fd unexpectedly succeeded")
	}
}
