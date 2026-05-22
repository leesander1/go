// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: AT THIS MOMENT REDOX CAN'T USE ALL OF THIS
// TODO: Should be merged with at_libc when openat is fully implemented

package unix

import (
	"syscall"
	"unsafe"
)

// Implemented as sysvicall6 in runtime/syscall_redox.go.
func syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// Implemented as rawsysvicall6 in runtime/syscall_redox.go.
func rawSyscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname procFstatat libc_fstatat

var (
	procFstatat uintptr
)

const (
	AT_FDCWD   = -0x64
	UTIME_OMIT = -0x2

	// not in relibc, but can be ignored by it
	AT_SYMLINK_NOFOLLOW = 0x100
	// not in relibc, and so unlinkat
	AT_REMOVEDIR = 0x200
)

func Unlinkat(dirfd int, path string, flags int) error {
	return syscall.ENOSYS
}

func Openat(dirfd int, path string, flags int, perm uint32) (int, error) {
	return -1, syscall.ENOSYS
}

func Fstatat(dirfd int, path string, stat *syscall.Stat_t, flags int) error {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return err
	}

	print("bout to sys call procFstatat\n")
	_, _, errno := syscall6(uintptr(unsafe.Pointer(&procFstatat)), 4,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(stat)),
		uintptr(flags),
		0, 0)
	if errno != 0 {
		return errno
	}

	return nil
}

func Readlinkat(dirfd int, path string, buf []byte) (int, error) {
	return -1, syscall.ENOSYS
}

func Mkdirat(dirfd int, path string, mode uint32) error {
	return syscall.ENOSYS
}

func Fchmodat(dirfd int, path string, mode uint32, flags int) error {
	return syscall.ENOSYS
}

func Fchownat(dirfd int, path string, uid, gid int, flags int) error {
	return syscall.ENOSYS
}

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) error {
	return syscall.ENOSYS
}

func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flag int) error {
	return syscall.ENOSYS
}

func Symlinkat(oldpath string, newdirfd int, newpath string) error {
	return syscall.ENOSYS
}

func Eaccess(path string, mode uint32) error {
	return syscall.ENOSYS
}
