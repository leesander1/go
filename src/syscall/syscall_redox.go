// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Redox system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_redox.go or syscall_unix.go.

package syscall

import "unsafe"

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

// We use cgo instead of asmcgo because of issues that ends in protection fault
// We cannot use cgocall in runtime because it uses malloc
func rawSysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

// linked by runtime.cgocall.go
//
//go:uintptrescapes
func cgocaller(unsafe.Pointer, ...uintptr) uintptr

// linked by runtime.cgocall.go
//
//go:uintptrescapes
func cgocaller2(unsafe.Pointer, ...uintptr) (r0 uintptr, err int32)

func syscgocall6(trap unsafe.Pointer, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	ret, errno := cgocaller2(trap, a1, a2, a3, a4, a5, a6)
	if errno != 0 {
		err = Errno(errno)
	} else {
		r1 = ret
	}
	return
}

func rawsyscgocall6(trap unsafe.Pointer, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	return syscgocall6(trap, nargs, a1, a2, a3, a4, a5, a6)
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Ino), unsafe.Sizeof(Dirent{}.Ino))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	reclen, ok := direntReclen(buf)
	if !ok {
		return 0, false
	}
	return reclen - uint64(unsafe.Offsetof(Dirent{}.Name)), true
}

func Pipe(p []int) (err error) {
	return Pipe2(p, 0)
}

//sysnb	pipe2(p *[2]_C_int, flags int) (err error)

func Pipe2(p []int, flags int) error {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err := pipe2(&pp, flags)
	if err == nil {
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	}
	return err
}

//TODO sys   accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) = libc.accept4

func Accept4(fd int, flags int) (int, Sockaddr, error) {
	panic("accept4 TODO")
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	// length is family (uint16), name, NUL.
	sl := _Socklen(2)
	if n > 0 {
		sl += _Socklen(n) + 1
	}
	if sa.raw.Path[0] == '@' || (sa.raw.Path[0] == 0 && sl > 3) {
		// Check sl > 3 so we don't change unnamed socket behavior.
		sa.raw.Path[0] = 0
		// Don't count trailing NUL for abstract address.
		sl--
	}

	return unsafe.Pointer(&sa.raw), sl, nil
}

func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(&rsa)
}

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, err error)

func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	// Getcwd will return an error if it failed for any reason.
	_, err = Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}

/*
 * Wrapped
 */

//sysnb	getgroups(ngid int, gid *_Gid_t) (n int, err error)
//sysnb	setgroups(ngid int, gid *_Gid_t) (err error)

func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	// Sanity check group count. Max is 16 on BSD.
	if n < 0 || n > 1000 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (err error) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

func ReadDirent(fd int, buf []byte) (n int, err error) {
	return PosixGetdents(fd, buf, 0)
}

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits.

type WaitStatus uint32

const (
	mask  = 0x7F
	core  = 0x80
	shift = 8

	exited  = 0
	stopped = 0x7F
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) ExitStatus() int {
	if w&mask != exited {
		return -1
	}
	return int(w >> shift)
}

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != 0 }

func (w WaitStatus) Signal() Signal {
	sig := Signal(w & mask)
	if sig == stopped || sig == 0 {
		return -1
	}
	return sig
}

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) Stopped() bool { return w&mask == stopped && Signal(w>>shift) != SIGSTOP }

func (w WaitStatus) Continued() bool { return w&mask == stopped && Signal(w>>shift) == SIGSTOP }

func (w WaitStatus) StopSignal() Signal {
	if !w.Stopped() {
		return -1
	}
	return Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int { return -1 }

func wait4(pid uintptr, wstatus *WaitStatus, options uintptr, rusage *Rusage) (wpid uintptr, err uintptr)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	r0, e1 := wait4(uintptr(pid), wstatus, uintptr(options), rusage)
	if e1 != 0 {
		err = Errno(e1)
	}
	return int(r0), err
}

func gethostname() (name string, err uintptr)

func Gethostname() (name string, err error) {
	name, e1 := gethostname()
	if e1 != 0 {
		err = Errno(e1)
	}
	return name, err
}

func UtimesNano(path string, ts []Timespec) error {
	panic("utimensat TODO")
}

//sys	fcntl(fd int, cmd int, arg int) (val int, err error)

// FcntlFlock performs a fcntl syscall for the [F_GETLK], [F_SETLK] or [F_SETLKW] command.
func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	panic("TODO FcntlFlock")
}

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		// Assume path ends at NUL.
		// This is not technically the redox semantics for
		// abstract Unix domain sockets -- they are supposed
		// to be uninterpreted fixed-size binary blobs -- but
		// everyone uses this convention.
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.Addr = pp.Addr
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		sa.Addr = pp.Addr
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error) = libc.accept

func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if err != nil {
		return
	}
	sa, err = anyToSockaddr(&rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func recvmsgRaw(fd int, p, oob []byte, flags int, rsa *RawSockaddrAny) (n, oobn int, recvflags int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(rsa))
	msg.Namelen = SizeofSockaddrAny
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// receive at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		// msg.Accrights = (*int8)(unsafe.Pointer(&oob[0]))
		// msg.Accrightslen = int32(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	// oobn = int(msg.Accrightslen)
	return
}

//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error) = libc.sendmsg

func sendmsgN(fd int, p, oob []byte, ptr unsafe.Pointer, salen _Socklen, flags int) (n int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = uint64(salen)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// send at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		// msg.Accrights = (*int8)(unsafe.Pointer(&oob[0]))
		// msg.Accrightslen = int32(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && len(p) == 0 {
		n = 0
	}
	return n, nil
}

/*
 * Exposed directly
 */
//sys	Access(path string, mode uint32) (err error)
//sys	Chdir(path string) (err error)
//sys	Chmod(path string, mode uint32) (err error)
//sys	Chown(path string, uid int, gid int) (err error)
//sys	Chroot(path string) (err error)
//sys	Close(fd int) (err error)
//sys	Dup(fd int) (nfd int, err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fpathconf(fd int, name int) (val int, err error)
//sys	Fstat(fd int, stat *Stat_t) (err error)
//TODO sys	Getdents(fd int, buf []byte, basep *uintptr) (n int, err error)
//sysnb	Getgid() (gid int)
//sysnb	Getpid() (pid int)
//sys	Geteuid() (euid int)
//sys	Getegid() (egid int)
//sys	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (n int, err error)
//sysnb	Getrlimit(which int, lim *Rlimit) (err error)
//sysnb	Getrusage(who int, rusage *Rusage) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)
//sysnb	Getuid() (uid int)
//sys	Kill(pid int, signum Signal) (err error)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Link(path string, link string) (err error)
//sys	Listen(s int, backlog int) (err error) = libc.listen
//sys	Lstat(path string, stat *Stat_t) (err error)
//sys	Mkdir(path string, mode uint32) (err error)
//sys	Mknod(path string, mode uint32, dev int) (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys	Open(path string, mode int, perm uint32) (fd int, err error)
//sys	Pathconf(path string, name int) (val int, err error)
//sys	pread(fd int, p []byte, offset int64) (n int, err error)
//sys	pwrite(fd int, p []byte, offset int64) (n int, err error)
//sys	read(fd int, p []byte) (n int, err error)
//sys	Readlink(path string, buf []byte) (n int, err error)
//sys	Rename(from string, to string) (err error)
//sys	Rmdir(path string) (err error)
//sys	Seek(fd int, offset int64, whence int) (newoffset int64, err error) = lseek
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) = libsendfile.sendfile
//sysnb	Setegid(egid int) (err error)
//sysnb	Seteuid(euid int) (err error)
//sysnb	Setgid(gid int) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sys	Setpriority(which int, who int, prio int) (err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sysnb	setrlimit(which int, lim *Rlimit) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Setuid(uid int) (err error)
//sys	Shutdown(s int, how int) (err error) = libc.shutdown
//sys	Stat(path string, stat *Stat_t) (err error)
//sys	Symlink(path string, link string) (err error)
//sys	Sync() (err error)
//sys	Truncate(path string, length int64) (err error)
//sys	PosixGetdents(fd int, buf []byte, flag int) (n int, err error)
//sys	Fsync(fd int) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sys	Umask(newmask int) (oldmask int)
//sys	Unlink(path string) (err error)
//sys	utimes(path string, times *[2]Timeval) (err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) = libc.bind
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) = libc.connect
//sys	mmap(addr uintptr, length uintptr, prot int, flag int, fd int, pos int64) (ret uintptr, err error)
//sys	munmap(addr uintptr, length uintptr) (err error)
//sys	sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error) = libc.sendto
//sys	socket(domain int, typ int, proto int) (fd int, err error) = libc.socket
//sysnb	socketpair(domain int, typ int, proto int, fd *[2]int32) (err error) = libc.socketpair
//sysnb	Uname(buf *Utsname) (err error) = libc.uname
//sys	write(fd int, p []byte) (n int, err error)
//sys	writev(fd int, iovecs []Iovec) (n uintptr, err error)
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error) = libc.getsockopt
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) = libc.getpeername
//sys	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) = libc.getsockname
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error) = libc.setsockopt
//sys	recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error) = libc.recvfrom
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error) = libc.recvmsg
//TODO sys	utimensat(dirfd int, path string, times *[2]Timespec, flag int) (err error)

func Getexecname() (path string, err error) {
	panic("Getexecname TODO")
}

func readlen(fd int, buf *byte, nbuf int) (n int, err error) {
	r0, _, e1 := syscgocall6(unsafe.Pointer(&libc_read), 3, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}

func Utimes(path string, tv []Timeval) error {
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

// Not implemented in Redox

type ICMPv6Filter struct {
	Filt [8]uint32
}

const SizeofICMPv6Filter = 0x20

const F_DUP2FD_CLOEXEC = 0x30
const _F_DUP2FD_CLOEXEC = F_DUP2FD_CLOEXEC
const F_DUPFD_CLOEXEC = 0
