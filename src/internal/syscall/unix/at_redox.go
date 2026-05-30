// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix

import (
	"runtime"
	"syscall"
	"unsafe"
)

// Implemented as sysvicall6 in runtime/syscall_redox.go.
func syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// Implemented as rawsysvicall6 in runtime/syscall_redox.go.
func rawSyscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:cgo_import_static _cgo_libc_openat
//go:cgo_import_static _cgo_libc_fstatat
//go:cgo_import_static _cgo_libc_readlinkat
//go:cgo_import_static _cgo_libc_mkdirat
//go:cgo_import_static _cgo_libc_fchmodat
//go:cgo_import_static _cgo_libc_fchownat
//go:cgo_import_static _cgo_libc_renameat
//go:cgo_import_static _cgo_libc_linkat
//go:cgo_import_static _cgo_libc_symlinkat
//go:cgo_import_static _cgo_libc_unlinkat

//go:linkname libc_openat _cgo_libc_openat
//go:linkname libc_fstatat _cgo_libc_fstatat
//go:linkname libc_readlinkat _cgo_libc_readlinkat
//go:linkname libc_mkdirat _cgo_libc_mkdirat
//go:linkname libc_fchmodat _cgo_libc_fchmodat
//go:linkname libc_fchownat _cgo_libc_fchownat
//go:linkname libc_renameat _cgo_libc_renameat
//go:linkname libc_linkat _cgo_libc_linkat
//go:linkname libc_symlinkat _cgo_libc_symlinkat
//go:linkname libc_unlinkat _cgo_libc_unlinkat

type libcFunc byte

var (
	libc_openat,
	libc_fstatat,
	libc_readlinkat,
	libc_mkdirat,
	libc_fchmodat,
	libc_fchownat,
	libc_renameat,
	libc_linkat,
	libc_symlinkat,
	libc_unlinkat libcFunc
)

const (
	AT_FDCWD   = -0x64
	UTIME_OMIT = -0x2

	AT_SYMLINK_NOFOLLOW = 0x200
	AT_REMOVEDIR        = 0x200
)

func Unlinkat(dirfd int, path string, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_unlinkat)), 3,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(flags),
		0, 0, 0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return errno
	}
	return nil
}

func Openat(dirfd int, path string, flags int, perm uint32) (int, error) {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}

	fd, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_openat)), 4,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(flags),
		uintptr(perm),
		0, 0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return 0, errno
	}
	return int(fd), nil
}

func Fstatat(dirfd int, path string, stat *syscall.Stat_t, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_fstatat)), 4,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(stat)),
		uintptr(flags),
		0, 0)
	runtime.KeepAlive(p)
	runtime.KeepAlive(stat)
	if errno != 0 {
		return errno
	}
	return nil
}

func Readlinkat(dirfd int, path string, buf []byte) (int, error) {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}
	var b *byte
	if len(buf) > 0 {
		b = &buf[0]
	}
	n, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_readlinkat)), 4,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(b)),
		uintptr(len(buf)),
		0, 0)
	runtime.KeepAlive(p)
	runtime.KeepAlive(buf)
	if errno != 0 {
		return 0, errno
	}
	return int(n), nil
}

func Mkdirat(dirfd int, path string, mode uint32) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_mkdirat)), 3,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(mode),
		0, 0, 0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return errno
	}
	return nil
}

func Fchmodat(dirfd int, path string, mode uint32, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_fchmodat)), 4,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(mode),
		uintptr(flags),
		0, 0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return errno
	}
	return nil
}

func Fchownat(dirfd int, path string, uid, gid int, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_fchownat)), 5,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(uid),
		uintptr(gid),
		uintptr(flags),
		0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return errno
	}
	return nil
}

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) error {
	oldp, err := syscall.BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := syscall.BytePtrFromString(newpath)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_renameat)), 4,
		uintptr(olddirfd),
		uintptr(unsafe.Pointer(oldp)),
		uintptr(newdirfd),
		uintptr(unsafe.Pointer(newp)),
		0, 0)
	runtime.KeepAlive(oldp)
	runtime.KeepAlive(newp)
	if errno != 0 {
		return errno
	}
	return nil
}

func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flag int) error {
	oldp, err := syscall.BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := syscall.BytePtrFromString(newpath)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_linkat)), 5,
		uintptr(olddirfd),
		uintptr(unsafe.Pointer(oldp)),
		uintptr(newdirfd),
		uintptr(unsafe.Pointer(newp)),
		uintptr(flag),
		0)
	runtime.KeepAlive(oldp)
	runtime.KeepAlive(newp)
	if errno != 0 {
		return errno
	}
	return nil
}

func Symlinkat(oldpath string, newdirfd int, newpath string) error {
	oldp, err := syscall.BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := syscall.BytePtrFromString(newpath)
	if err != nil {
		return err
	}

	_, _, errno := syscall6(uintptr(unsafe.Pointer(&libc_symlinkat)), 3,
		uintptr(unsafe.Pointer(oldp)),
		uintptr(newdirfd),
		uintptr(unsafe.Pointer(newp)),
		0, 0, 0)
	runtime.KeepAlive(oldp)
	runtime.KeepAlive(newp)
	if errno != 0 {
		return errno
	}
	return nil
}

func Eaccess(path string, mode uint32) error {
	return syscall.ENOSYS
}
