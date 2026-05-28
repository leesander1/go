// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package cache

import "fmt"

func cacheFileName(id [HashSize]byte, key string) (dir, name string) {
	return fmt.Sprintf("%02x", id[0]), fmt.Sprintf("%x", id) + "-" + key
}
