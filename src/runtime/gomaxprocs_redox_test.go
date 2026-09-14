// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package runtime_test

import (
	"runtime"
	"testing"
)

func TestRedoxGOMAXPROCSStaysSingleP(t *testing.T) {
	if old := runtime.GOMAXPROCS(8); old != 1 {
		t.Fatalf("GOMAXPROCS(8) returned %d, want 1", old)
	}
	if got := runtime.GOMAXPROCS(0); got != 1 {
		t.Fatalf("GOMAXPROCS(0) = %d, want 1", got)
	}
}
