// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package unix_test

import (
	"errors"
	"internal/syscall/unix"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestEaccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program")
	if err := os.WriteFile(path, []byte("program"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := unix.Eaccess(path, unix.X_OK); !errors.Is(err, syscall.EACCES) {
		t.Fatalf("Eaccess(non-executable) = %v; want EACCES", err)
	}
	if err := os.Chmod(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := unix.Eaccess(path, unix.X_OK); err != nil {
		t.Fatalf("Eaccess(executable) = %v; want nil", err)
	}
	if err := unix.Eaccess(path+"-missing", 0); !errors.Is(err, syscall.ENOENT) {
		t.Fatalf("Eaccess(missing) = %v; want ENOENT", err)
	}
}
