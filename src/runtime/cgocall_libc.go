// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package runtime

import "unsafe"

func cgocallLibc(fn, arg unsafe.Pointer) {
	cgocall(fn, arg)
}

//go:nosplit
func asmcgocallLibc(fn, arg unsafe.Pointer) {
	asmcgocall(fn, arg)
}
