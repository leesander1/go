// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package cache

import (
	"os"
	"time"
)

func cacheChtimes(name string, atime, mtime time.Time) {
	// Cache mtimes are approximate metadata for trimming. Avoid Redox
	// utimensat here until it is reliable under the Go build workload.
}

func cacheIndexOpenMode() int {
	return os.O_RDWR | os.O_CREATE
}
