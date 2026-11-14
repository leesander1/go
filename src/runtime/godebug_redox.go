// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/stringslite"
	"unsafe"
)

func getGodebugEarlyRedox(prefix string) string {
	// Redox exposes the environment through relibc. Unlike other Unix-like
	// systems, its environment does not immediately follow argv.
	envp := (**byte)(unsafe.Pointer(get_environ()))
	for i := int32(0); ; i++ {
		p := argv_index(envp, i)
		if p == nil {
			return ""
		}
		s := unsafe.String(p, findnull(p))
		if stringslite.HasPrefix(s, prefix) {
			return gostringnocopy(p)[len(prefix):]
		}
	}
}
