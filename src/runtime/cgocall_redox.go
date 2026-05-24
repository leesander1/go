// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package runtime

import "unsafe"

//go:linkname syscall_rawcgocaller2 syscall.rawcgocaller2
//go:nosplit
//go:uintptrescapes
func syscall_rawcgocaller2(fn unsafe.Pointer, args ...uintptr) (r0 uintptr, err int32) {
	as := argset{args: unsafe.Pointer(&args[0])}
	asmcgocall(fn, unsafe.Pointer(&as))
	r0 = as.retval
	err = as.errno
	return
}
