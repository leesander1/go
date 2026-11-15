// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall_test

import (
	"bytes"
	"fmt"
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

func TestMmapFileOffsetContents(t *testing.T) {
	pageSize := syscall.Getpagesize()
	want := make([]byte, 12*pageSize)
	for i := range want {
		want[i] = byte(i*37 + i/pageSize)
	}
	name := t.TempDir() + "/contents"
	if err := os.WriteFile(name, want, 0600); err != nil {
		t.Fatal(err)
	}
	for _, flags := range []int{syscall.MAP_SHARED, syscall.MAP_PRIVATE} {
		for _, page := range []int{0, 1, 3} {
			t.Run(fmt.Sprintf("flags%d/page%d", flags, page), func(t *testing.T) {
				f, err := os.Open(name)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				offset, length := page*pageSize, 5*pageSize+37
				b, err := syscall.Mmap(int(f.Fd()), int64(offset), length, syscall.PROT_READ, flags)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := syscall.Munmap(b); err != nil {
						t.Error(err)
					}
				}()
				check := func(phase string) {
					t.Helper()
					for i, got := range b {
						if got != want[offset+i] {
							t.Fatalf("%s: byte %d = %#x, want %#x", phase, i, got, want[offset+i])
						}
					}
				}
				check("open file")
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
				check("closed file")
			})
		}
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
