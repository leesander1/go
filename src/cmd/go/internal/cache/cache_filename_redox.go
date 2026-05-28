// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package cache

import "encoding/hex"

func cacheFileName(id [HashSize]byte, key string) (dir, name string) {
	var encoded [HashSize * 2]byte
	hex.Encode(encoded[:], id[:])
	return string(encoded[:2]), string(encoded[:]) + "-" + key
}
