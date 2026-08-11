// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package runtime

import (
	"internal/abi"
	"unsafe"
)

type redoxCgoCallArgs struct {
	fn  unsafe.Pointer
	arg unsafe.Pointer
}

func redoxCgoCall()

func cgocallLibc(fn, arg unsafe.Pointer) {
	call := redoxCgoCallArgs{fn: fn, arg: arg}
	cgocall(unsafe.Pointer(abi.FuncPCABI0(redoxCgoCall)), unsafe.Pointer(&call))
}

//go:nosplit
func asmcgocallLibc(fn, arg unsafe.Pointer) {
	call := redoxCgoCallArgs{fn: fn, arg: arg}
	asmcgocall(unsafe.Pointer(abi.FuncPCABI0(redoxCgoCall)), unsafe.Pointer(&call))
}

//go:linkname syscall_cgocaller6 syscall.cgocaller6
//go:uintptrescapes
func syscall_cgocaller6(fn unsafe.Pointer, a1, a2, a3, a4, a5, a6 uintptr) (r0 uintptr, err int32) {
	args := [6]uintptr{a1, a2, a3, a4, a5, a6}
	as := argset{args: unsafe.Pointer(&args[0])}
	entersyscallblock()
	asmcgocallLibc(fn, unsafe.Pointer(&as))
	exitsyscall()
	r0 = as.retval
	err = as.errno
	return
}

//go:linkname syscall_rawcgocaller2 syscall.rawcgocaller2
//go:nosplit
//go:uintptrescapes
func syscall_rawcgocaller2(fn unsafe.Pointer, args ...uintptr) (r0 uintptr, err int32) {
	as := argset{args: unsafe.Pointer(&args[0])}
	asmcgocallLibc(fn, unsafe.Pointer(&as))
	r0 = as.retval
	err = as.errno
	return
}

//go:linkname syscall_rawcgocaller6 syscall.rawcgocaller6
//go:nosplit
//go:uintptrescapes
func syscall_rawcgocaller6(fn unsafe.Pointer, a1, a2, a3, a4, a5, a6 uintptr) (r0 uintptr, err int32) {
	args := [6]uintptr{a1, a2, a3, a4, a5, a6}
	as := argset{args: unsafe.Pointer(&args[0])}
	asmcgocallLibc(fn, unsafe.Pointer(&as))
	r0 = as.retval
	err = as.errno
	return
}
