// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMapFileUnalignedOffset(t *testing.T) {
	const contents = "prefix: mapped compiler import data :suffix"
	path := filepath.Join(t.TempDir(), "mapfile")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	const want = "mapped compiler import data"
	got, err := MapFile(f, int64(len("prefix: ")), int64(len(want)))
	if err != nil {
		t.Fatalf("MapFile: %v", err)
	}
	if got != want {
		t.Fatalf("MapFile = %q; want %q", got, want)
	}
}
