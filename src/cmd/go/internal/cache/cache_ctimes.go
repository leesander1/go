// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package cache

import (
	"os"
	"time"
)

func cacheChtimes(name string, atime, mtime time.Time) {
	os.Chtimes(name, atime, mtime)
}

func cacheIndexOpenMode() int {
	return os.O_WRONLY | os.O_CREATE
}
