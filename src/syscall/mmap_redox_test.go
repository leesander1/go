// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall_test

import (
	"bytes"
	"os"
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

func TestMmapSharedFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "mmap")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	pageSize := syscall.Getpagesize()
	if err := f.Truncate(int64(pageSize)); err != nil {
		t.Fatal(err)
	}
	b, err := syscall.Mmap(int(f.Fd()), 0, pageSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FILE)
	if err != nil {
		t.Fatalf("Mmap: %v", err)
	}
	want := []byte("redox file-backed mmap")
	copy(b, want)
	if err := syscall.Munmap(b); err != nil {
		t.Fatalf("Munmap: %v", err)
	}

	got := make([]byte, len(want))
	if _, err := f.ReadAt(got, 0); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("mapped write = %q; want %q", got, want)
	}
}
