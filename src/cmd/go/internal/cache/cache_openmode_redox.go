// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package cache

import "os"

func cacheIndexOpenMode() int {
	// Redox can stall opening an index for reading after it was created
	// through a write-only descriptor. Keep the descriptor read/write.
	return os.O_RDWR | os.O_CREATE
}
