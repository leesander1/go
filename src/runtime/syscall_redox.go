// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

//go:cgo_import_static _cgo_libc_chdir
//go:cgo_import_static _cgo_libc_chroot
//go:cgo_import_static _cgo_libc_close
//go:cgo_import_static _cgo_libc_dup2
//go:cgo_import_static _cgo_libc_execve
//go:cgo_import_static _cgo_libc_fork
//go:cgo_import_static _cgo_libc_fcntl
//go:cgo_import_static _cgo_libc_gethostname
//go:cgo_import_static _cgo_libc_getpid
//go:cgo_import_static _cgo_libc_ioctl
//go:cgo_import_static _cgo_libc_setgid
//go:cgo_import_static _cgo_libc_setgroups
//go:cgo_import_static _cgo_libc_setrlimit
//go:cgo_import_static _cgo_libc_setsid
//go:cgo_import_static _cgo_libc_setuid
//go:cgo_import_static _cgo_libc_setpgid
//go:cgo_import_static _cgo_libc_issetugid

//go:linkname libc_chdir _cgo_libc_chdir
//go:linkname libc_chroot _cgo_libc_chroot
//go:linkname libc_close _cgo_libc_close
//go:linkname libc_dup2 _cgo_libc_dup2
//go:linkname libc_execve _cgo_libc_execve
//go:linkname libc_fcntl _cgo_libc_fcntl
//go:linkname libc_fork _cgo_libc_fork
//go:linkname libc_gethostname _cgo_libc_gethostname
//go:linkname libc_getpid _cgo_libc_getpid
//go:linkname libc_ioctl _cgo_libc_ioctl
//go:linkname libc_setgid _cgo_libc_setgid
//go:linkname libc_setgroups _cgo_libc_setgroups
//go:linkname libc_setrlimit _cgo_libc_setrlimit
//go:linkname libc_setsid _cgo_libc_setsid
//go:linkname libc_setuid _cgo_libc_setuid
//go:linkname libc_setpgid _cgo_libc_setpgid
//go:linkname libc_issetugid _cgo_libc_issetugid

var (
	libc_chdir,
	libc_chroot,
	libc_close,
	libc_dup2,
	libc_execve,
	libc_fork,
	libc_fcntl,
	libc_gethostname,
	libc_getpid,
	libc_ioctl,
	libc_setgid,
	libc_setgroups,
	libc_setrlimit,
	libc_setsid,
	libc_setuid,
	libc_setpgid,
	libc_issetugid byte
)

// Many of these are exported via linkname to assembly in the syscall
// package.

//go:nosplit
//go:linkname syscall_sysvicall6
//go:cgo_unsafe_args
func syscall_sysvicall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	entersyscallblock()
	r1, r2, err = syscall_rawsysvicall6(fn, nargs, a1, a2, a3, a4, a5, a6)
	exitsyscall()
	return
}

//go:nosplit
//go:linkname syscall_rawsysvicall6
//go:cgo_unsafe_args
func syscall_rawsysvicall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	args := [6]uintptr{a1, a2, a3, a4, a5, a6}
	call := libcall{
		fn:   fn,
		n:    nargs,
		args: uintptr(noescape(unsafe.Pointer(&args[0]))),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, 0, call.err
}

// TODO(aram): Once we remove all instances of C calling sysvicallN, make
// sysvicallN return errors and replace the body of the following functions
// with calls to sysvicallN.

//go:nosplit
//go:linkname syscall_chdir
func syscall_chdir(path uintptr) (err uintptr) {
	_, errno := cgocaller1(unsafe.Pointer(&libc_chdir), path)
	return uintptr(errno)
}

//go:nosplit
//go:linkname syscall_chroot
func syscall_chroot(path uintptr) (err uintptr) {
	_, errno := cgocaller1(unsafe.Pointer(&libc_chroot), path)
	return uintptr(errno)
}

// like close, but must not split stack, for forkx.
//
//go:nosplit
//go:linkname syscall_close
func syscall_close(fd int32) int32 {
	_, errno := cgocaller1(unsafe.Pointer(&libc_close), uintptr(fd))
	return errno
}

//go:nosplit
//go:linkname syscall_dup2
func syscall_dup2(oldfd, newfd uintptr) (val, err uintptr) {
	ret, errno := cgocaller2(unsafe.Pointer(&libc_dup2), oldfd, newfd)
	if errno != 0 {
		err = uintptr(errno)
	} else {
		val = ret
	}
	return
}

//go:nosplit
//go:linkname syscall_execve
//go:cgo_unsafe_args
func syscall_execve(path, argv, envp uintptr) (err uintptr) {
	_, errno := cgocaller3(unsafe.Pointer(&libc_execve), path, argv, envp)
	return uintptr(errno)
}

//go:nosplit
func syscall_execve_stack(path, argv, envp, stack uintptr) (err uintptr)

// like exit, but must not split stack, for forkx.
//
//go:nosplit
//go:linkname syscall_exit
func syscall_exit(code uintptr) {
	cgocaller1(unsafe.Pointer(&libc_exit), code)
}

//go:nosplit
//go:linkname syscall_fcntl
//go:cgo_unsafe_args
func syscall_fcntl(fd, cmd, arg uintptr) (val, err uintptr) {
	ret, errno := cgocaller3(unsafe.Pointer(&libc_fcntl), fd, cmd, arg)
	if errno != 0 {
		err = uintptr(errno)
	} else {
		val = ret
	}
	return
}

//go:nosplit
//go:linkname syscall_forkx
func syscall_forkx(_flags uintptr) (pid uintptr, err uintptr) {
	ret, errno := cgocaller0(unsafe.Pointer(&libc_fork))
	if errno != 0 {
		err = uintptr(errno)
	} else {
		pid = ret
	}
	return
}

//go:linkname syscall_gethostname
func syscall_gethostname() (name string, err uintptr) {
	cname := new([_MAXHOSTNAMELEN]byte)
	_, errno := cgocaller2(unsafe.Pointer(&libc_gethostname), uintptr(unsafe.Pointer(&cname[0])), _MAXHOSTNAMELEN)
	if errno != 0 {
		return "", uintptr(errno)
	}
	cname[_MAXHOSTNAMELEN-1] = 0
	return gostringnocopy(&cname[0]), 0
}

//go:nosplit
//go:linkname syscall_getpid
func syscall_getpid() (pid, err uintptr) {
	ret, errno := cgocaller0(unsafe.Pointer(&libc_getpid))
	return ret, uintptr(errno)
}

//go:nosplit
//go:linkname syscall_ioctl
//go:cgo_unsafe_args
func syscall_ioctl(fd, req, arg uintptr) (err uintptr) {
	_, errno := cgocaller3(unsafe.Pointer(&libc_ioctl), fd, req, arg)
	return uintptr(errno)
}

// This is syscall.RawSyscall, it exists to satisfy some build dependency,
// but it doesn't work.
//
//go:linkname syscall_rawsyscall
func syscall_rawsyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	panic("RawSyscall not available on Redox")
}

// This is syscall.RawSyscall6, it exists to avoid a linker error because
// syscall.RawSyscall6 is already declared. See golang.org/issue/24357
//
//go:linkname syscall_rawsyscall6
func syscall_rawsyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	panic("RawSyscall6 not available on Redox")
}

//go:nosplit
//go:linkname syscall_setgid
func syscall_setgid(gid uintptr) (err uintptr) {
	_, errno := cgocaller1(unsafe.Pointer(&libc_setgid), gid)
	return uintptr(errno)
}

//go:nosplit
//go:linkname syscall_setgroups
//go:cgo_unsafe_args
func syscall_setgroups(ngid, gid uintptr) (err uintptr) {
	_, errno := cgocaller2(unsafe.Pointer(&libc_setgroups), ngid, gid)
	return uintptr(errno)
}

func issetugid() int32 {
	_, errno := cgocaller0(unsafe.Pointer(&libc_issetugid))
	return errno
}

//go:nosplit
//go:linkname syscall_setrlimit
//go:cgo_unsafe_args
func syscall_setrlimit(which uintptr, lim unsafe.Pointer) (err uintptr) {
	_, errno := cgocaller2(unsafe.Pointer(&libc_setrlimit), which, uintptr(lim))
	KeepAlive(lim)
	return uintptr(errno)
}

//go:nosplit
//go:linkname syscall_setsid
func syscall_setsid() (pid, err uintptr) {
	ret, errno := cgocaller0(unsafe.Pointer(&libc_setsid))
	return ret, uintptr(errno)
}

//go:nosplit
//go:linkname syscall_setuid
func syscall_setuid(uid uintptr) (err uintptr) {
	_, errno := cgocaller1(unsafe.Pointer(&libc_setuid), uid)
	return uintptr(errno)
}

//go:nosplit
//go:linkname syscall_setpgid
//go:cgo_unsafe_args
func syscall_setpgid(pid, pgid uintptr) (err uintptr) {
	_, errno := cgocaller2(unsafe.Pointer(&libc_setpgid), pid, pgid)
	return uintptr(errno)
}

//go:linkname syscall_syscall
//go:cgo_unsafe_args
func syscall_syscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	panic("no syscall on Redox")
}

//go:linkname syscall_wait4
//go:cgo_unsafe_args
func syscall_wait4(pid uintptr, wstatus *uint32, options uintptr, rusage unsafe.Pointer) (wpid int, err uintptr) {
	const _WNOHANG = 0x1

	waitOptions := options
	blocking := options&_WNOHANG == 0
	if blocking {
		waitOptions |= _WNOHANG
	}

	for {
		ret, errno := cgocaller3(unsafe.Pointer(&libc_waitpid), pid, uintptr(unsafe.Pointer(wstatus)), waitOptions)
		KeepAlive(wstatus)
		KeepAlive(rusage)
		if errno != 0 {
			err = uintptr(errno)
			return
		}
		if ret != 0 || !blocking {
			wpid = int(ret)
			return
		}

		// Redox waitpid goes through libc/proc-scheme code. Polling with
		// WNOHANG keeps the wait cooperative so os/exec copy goroutines can
		// run while the parent is waiting for the child to exit.
		Gosched()
		entersyscallblock()
		usleep_no_g(1000)
		exitsyscall()
	}
}

//go:nosplit
//go:linkname syscall_write
//go:cgo_unsafe_args
func syscall_write(fd, buf, nbyte uintptr) (n, err uintptr) {
	ret, errno := cgocaller3(unsafe.Pointer(&libc_write), fd, buf, nbyte)
	if errno != 0 {
		err = uintptr(errno)
	} else {
		n = ret
	}
	return
}
