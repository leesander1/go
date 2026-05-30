// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox && amd64

package syscall

import (
	"runtime"
	"unsafe"
)

//go:cgo_import_static _cgo_libc_utimensat
//go:linkname libc_utimensat _cgo_libc_utimensat

var libc_utimensat libcFunc

func utimensat(dirfd int, path string, times *[2]Timespec, flag int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := syscgocall6(unsafe.Pointer(&libc_utimensat), 4, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(times)), uintptr(flag), 0, 0)
	runtime.KeepAlive(_p0)
	runtime.KeepAlive(times)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
