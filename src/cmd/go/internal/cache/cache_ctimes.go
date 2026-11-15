// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"os"
	"time"
)

func cacheChtimes(name string, atime, mtime time.Time) {
	os.Chtimes(name, atime, mtime)
}
