// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/runtime/atomic"
	"unsafe"
)

//go:cgo_export_dynamic runtime.end _end
//go:cgo_export_dynamic runtime.etext _etext
//go:cgo_export_dynamic runtime.edata _edata

//go:cgo_import_dynamic libc___errno __errno "libc.so"
//go:cgo_import_static _cgo_libc_environ
//go:cgo_import_dynamic libc_clock_gettime clock_gettime "libc.so"
//go:cgo_import_static _cgo_libc_exit
//go:cgo_import_static _cgo_libc_waitpid
//go:cgo_import_dynamic libc_getcontext getcontext "libc.so"
//go:cgo_import_dynamic libc_kill kill "libc.so"
//go:cgo_import_dynamic libc_madvise madvise "libc.so"
//go:cgo_import_dynamic libc_malloc malloc "libc.so"
//go:cgo_import_dynamic libc_mmap mmap "libc.so"
//go:cgo_import_dynamic libc_mprotect mprotect "libc.so"
//go:cgo_import_dynamic libc_munmap munmap "libc.so"
//go:cgo_import_dynamic libc_open open "libc.so"
//go:cgo_import_dynamic libc_pthread_attr_destroy pthread_attr_destroy "libc.so"
//go:cgo_import_dynamic libc_pthread_attr_getstack pthread_attr_getstack "libc.so"
//go:cgo_import_dynamic libc_pthread_attr_init pthread_attr_init "libc.so"
//go:cgo_import_dynamic libc_pthread_attr_setdetachstate pthread_attr_setdetachstate "libc.so"
//go:cgo_import_dynamic libc_pthread_attr_setstack pthread_attr_setstack "libc.so"
//go:cgo_import_dynamic libc_pthread_create pthread_create "libc.so"
//go:cgo_import_dynamic libc_pthread_getattr_np pthread_getattr_np "libc.so"
//go:cgo_import_dynamic libc_pthread_self pthread_self "libc.so"
//go:cgo_import_dynamic libc_pthread_kill pthread_kill "libc.so"
//go:cgo_import_dynamic libc_raise raise "libc.so"
//go:cgo_import_dynamic libc_read read "libc.so"
//go:cgo_import_dynamic libc_select select "libc.so"
//go:cgo_import_dynamic libc_sched_yield sched_yield "libc.so"
//go:cgo_import_dynamic libc_sem_init sem_init "libc.so"
//go:cgo_import_dynamic libc_sem_post sem_post "libc.so"
//go:cgo_import_dynamic libc_sem_timedwait sem_timedwait "libc.so"
//go:cgo_import_dynamic libc_sem_wait sem_wait "libc.so"
//go:cgo_import_dynamic libc_setitimer setitimer "libc.so"
//go:cgo_import_dynamic libc_sigaction sigaction "libc.so"
//go:cgo_import_dynamic libc_sigaltstack sigaltstack "libc.so"
//go:cgo_import_dynamic libc_sigprocmask sigprocmask "libc.so"
//go:cgo_import_dynamic libc_sysconf sysconf "libc.so"
//go:cgo_import_dynamic libc_usleep usleep "libc.so"
//go:cgo_import_dynamic libc_write write "libc.so"
//go:cgo_import_dynamic libc_pipe2 pipe2 "libc.so"

//go:linkname libc___errno libc___errno
//go:linkname libc_environ _cgo_libc_environ
//go:linkname libc_clock_gettime libc_clock_gettime
//go:linkname libc_exit _cgo_libc_exit
//go:linkname libc_waitpid _cgo_libc_waitpid
//go:linkname libc_getcontext libc_getcontext
//go:linkname libc_kill libc_kill
//go:linkname libc_madvise libc_madvise
//go:linkname libc_malloc libc_malloc
//go:linkname libc_mmap libc_mmap
//go:linkname libc_mprotect libc_mprotect
//go:linkname libc_munmap libc_munmap
//go:linkname libc_open libc_open
//go:linkname libc_pthread_attr_destroy libc_pthread_attr_destroy
//go:linkname libc_pthread_attr_getstack libc_pthread_attr_getstack
//go:linkname libc_pthread_attr_init libc_pthread_attr_init
//go:linkname libc_pthread_attr_setdetachstate libc_pthread_attr_setdetachstate
//go:linkname libc_pthread_attr_setstack libc_pthread_attr_setstack
//go:linkname libc_pthread_create libc_pthread_create
//go:linkname libc_pthread_getattr_np libc_pthread_getattr_np
//go:linkname libc_pthread_self libc_pthread_self
//go:linkname libc_pthread_kill libc_pthread_kill
//go:linkname libc_raise libc_raise
//go:linkname libc_read libc_read
//go:linkname libc_select libc_select
//go:linkname libc_sched_yield libc_sched_yield
//go:linkname libc_sem_init libc_sem_init
//go:linkname libc_sem_post libc_sem_post
//go:linkname libc_sem_timedwait libc_sem_timedwait
//go:linkname libc_sem_wait libc_sem_wait
//go:linkname libc_setitimer libc_setitimer
//go:linkname libc_sigaction libc_sigaction
//go:linkname libc_sigaltstack libc_sigaltstack
//go:linkname libc_sigprocmask libc_sigprocmask
//go:linkname libc_sysconf libc_sysconf
//go:linkname libc_usleep libc_usleep
//go:linkname libc_write libc_write
//go:linkname libc_pipe2 libc_pipe2

var (
	libc___errno,
	libc_clock_gettime,
	libc_getcontext,
	libc_kill,
	libc_madvise,
	libc_malloc,
	libc_mmap,
	libc_mprotect,
	libc_munmap,
	libc_open,
	libc_pthread_attr_destroy,
	libc_pthread_attr_getstack,
	libc_pthread_attr_init,
	libc_pthread_attr_setdetachstate,
	libc_pthread_attr_setstack,
	libc_pthread_create,
	libc_pthread_getattr_np,
	libc_pthread_self,
	libc_pthread_kill,
	libc_raise,
	libc_read,
	libc_sched_yield,
	libc_select,
	libc_sem_init,
	libc_sem_post,
	libc_sem_timedwait,
	libc_sem_wait,
	libc_setitimer,
	libc_sigaction,
	libc_sigaltstack,
	libc_sigprocmask,
	libc_sysconf,
	libc_usleep,
	libc_write,
	libc_pipe2 libcFunc
)

var (
	libc_environ,
	libc_exit,
	libc_waitpid byte
)

//go:nosplit
func closefd(fd int32) int32 {
	_, errno := cgocaller1(unsafe.Pointer(&libc_close), uintptr(fd))
	return errno
}

//go:nosplit
func exit(r int32) {
	cgocaller1(unsafe.Pointer(&libc_exit), uintptr(r))
}

//go:nosplit
func get_environ() uintptr {
	ret, _ := cgocaller1(unsafe.Pointer(&libc_environ), uintptr(0))
	return ret
}

//go:nosplit
func getcontext(context *ucontext) /* int32 */ {
	sysvicall1(&libc_getcontext, uintptr(unsafe.Pointer(context)))
	KeepAlive(context)
}

//go:nosplit
func madvise(addr unsafe.Pointer, n uintptr, flags int32) {
	sysvicall3(&libc_madvise, uintptr(addr), uintptr(n), uintptr(flags))
	KeepAlive(addr)
}

//go:nosplit
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (unsafe.Pointer, int) {
	p, err := doMmap(uintptr(addr), n, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(off))
	KeepAlive(addr)
	if p == ^uintptr(0) {
		return nil, int(err)
	}
	return unsafe.Pointer(p), 0
}

//go:nosplit
//go:cgo_unsafe_args
func doMmap(addr, n, prot, flags, fd, off uintptr) (uintptr, uintptr) {
	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(&libc_mmap))
	libcall.n = 6
	libcall.args = uintptr(noescape(unsafe.Pointer(&addr)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	return libcall.r1, libcall.err
}

//go:nosplit
func munmap(addr unsafe.Pointer, n uintptr) {
	sysvicall2(&libc_munmap, uintptr(addr), uintptr(n))
	KeepAlive(addr)
}

//go:nosplit
func mprotect(addr unsafe.Pointer, n uintptr, prot int32) int32 {
	r1, err := sysvicall3Err(&libc_mprotect, uintptr(addr), uintptr(n), uintptr(prot))
	KeepAlive(addr)
	if r1 == ^uintptr(0) {
		return int32(err)
	}
	return 0
}

//go:nosplit
func nanotime1() int64 {
	var ts timespec
	sysvicall2(&libc_clock_gettime, CLOCK_MONOTONIC, uintptr(unsafe.Pointer(&ts)))
	return ts.tv_sec*1e9 + ts.tv_nsec
}

//go:nosplit
func open(path *byte, mode, perm int32) int32 {
	ret := int32(sysvicall3(&libc_open, uintptr(unsafe.Pointer(path)), uintptr(mode), uintptr(perm)))
	KeepAlive(path)
	return ret
}

func pthread_attr_destroy(attr *pthread_attr_t) int32 {
	ret := int32(sysvicall1(&libc_pthread_attr_destroy, uintptr(unsafe.Pointer(attr))))
	KeepAlive(attr)
	return ret
}

func pthread_attr_getstack(attr *pthread_attr_t, addr unsafe.Pointer, size *uint64) int32 {
	ret := int32(sysvicall3(&libc_pthread_attr_getstack, uintptr(unsafe.Pointer(attr)), uintptr(addr), uintptr(unsafe.Pointer(size))))
	KeepAlive(attr)
	KeepAlive(addr)
	KeepAlive(size)
	return ret
}

func pthread_attr_init(attr *pthread_attr_t) int32 {
	ret := int32(sysvicall1(&libc_pthread_attr_init, uintptr(unsafe.Pointer(attr))))
	KeepAlive(attr)
	return ret
}

func pthread_attr_setdetachstate(attr *pthread_attr_t, state int32) int32 {
	ret := int32(sysvicall2(&libc_pthread_attr_setdetachstate, uintptr(unsafe.Pointer(attr)), uintptr(state)))
	KeepAlive(attr)
	return ret
}

func pthread_attr_setstack(attr *pthread_attr_t, addr uintptr, size uint64) int32 {
	ret := int32(sysvicall3(&libc_pthread_attr_setstack, uintptr(unsafe.Pointer(attr)), uintptr(addr), uintptr(size)))
	KeepAlive(attr)
	return ret
}

func pthread_create(thread *pthread_t, attr *pthread_attr_t, fn uintptr, arg unsafe.Pointer) int32 {
	ret := int32(sysvicall4(&libc_pthread_create, uintptr(unsafe.Pointer(thread)), uintptr(unsafe.Pointer(attr)), uintptr(fn), uintptr(arg)))
	KeepAlive(thread)
	KeepAlive(attr)
	KeepAlive(arg)
	return ret
}

func pthread_getattr_np(thread pthread_t, attr *pthread_attr_t) int32 {
	ret := int32(sysvicall2(&libc_pthread_getattr_np, uintptr(thread), uintptr(unsafe.Pointer(attr))))
	KeepAlive(attr)
	return ret
}

func pthread_self() pthread_t {
	return pthread_t(sysvicall0(&libc_pthread_self))
}

func signalM(mp *m, sig int) {
	sysvicall2(&libc_pthread_kill, uintptr(pthread_t(mp.procid)), uintptr(sig))
}

//go:nosplit
//go:nowritebarrierrec
func raise(sig uint32) /* int32 */ {
	sysvicall1(&libc_raise, uintptr(sig))
}

func raiseproc(sig uint32) /* int32 */ {
	pid, _ := syscall_getpid()
	sysvicall2(&libc_kill, pid, uintptr(sig))
}

//go:nosplit
func read(fd int32, buf unsafe.Pointer, nbyte int32) int32 {
	r1, err := sysvicall3Err(&libc_read, uintptr(fd), uintptr(buf), uintptr(nbyte))
	KeepAlive(buf)
	if c := int32(r1); c >= 0 {
		return c
	}
	return -int32(err)
}

//go:nosplit
func sem_init(sem *sem_t, pshared int32, value uint32) int32 {
	ret := int32(sysvicall3(&libc_sem_init, uintptr(unsafe.Pointer(sem)), uintptr(pshared), uintptr(value)))
	KeepAlive(sem)
	return ret
}

//go:nosplit
func sem_post(sem *sem_t) int32 {
	ret := int32(sysvicall1(&libc_sem_post, uintptr(unsafe.Pointer(sem))))
	KeepAlive(sem)
	return ret
}

//go:nosplit
func sem_timedwait(sem *sem_t, timeout *timespec) int32 {
	ret := int32(sysvicall2(&libc_sem_timedwait, uintptr(unsafe.Pointer(sem)), uintptr(unsafe.Pointer(timeout))))
	KeepAlive(sem)
	KeepAlive(timeout)
	return ret
}

//go:nosplit
func sem_wait(sem *sem_t) int32 {
	ret := int32(sysvicall1(&libc_sem_wait, uintptr(unsafe.Pointer(sem))))
	KeepAlive(sem)
	return ret
}

func setitimer(which int32, value *itimerval, ovalue *itimerval) /* int32 */ {
	sysvicall3(&libc_setitimer, uintptr(which), uintptr(unsafe.Pointer(value)), uintptr(unsafe.Pointer(ovalue)))
	KeepAlive(value)
	KeepAlive(ovalue)
}

//go:nosplit
//go:nowritebarrierrec
func sigaction(sig uint32, act *sigactiont, oact *sigactiont) /* int32 */ {
	sysvicall3(&libc_sigaction, uintptr(sig), uintptr(unsafe.Pointer(act)), uintptr(unsafe.Pointer(oact)))
	KeepAlive(act)
	KeepAlive(oact)
}

//go:nosplit
//go:nowritebarrierrec
func sigaltstack(ss *stackt, oss *stackt) /* int32 */ {
	sysvicall2(&libc_sigaltstack, uintptr(unsafe.Pointer(ss)), uintptr(unsafe.Pointer(oss)))
	KeepAlive(ss)
	KeepAlive(oss)
}

//go:nosplit
//go:nowritebarrierrec
func sigprocmask(how int32, set *sigset, oset *sigset) /* int32 */ {
	sysvicall3(&libc_sigprocmask, uintptr(how), uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(oset)))
	KeepAlive(set)
	KeepAlive(oset)
}

//go:nosplit
//go:nowritebarrierrec
func forkRestoreSigmask(sigmask sigset) {
	// The post-fork path cannot safely call back into relibc's signal
	// machinery, which may hold locks copied from other threads.
	redoxSigprocmask(_SIG_SETMASK, &sigmask, nil)
}

type redoxSigcontrol struct {
	word            [2]uint64
	senderInfos     [32]uint64
	controlFlags    uintptr
	savedIP         uintptr
	savedArchdepReg uintptr
}

func redoxSigcontrolCurrent() *redoxSigcontrol

const redoxSigmaskAllowAlways = uint64(1<<(_SIGKILL-1) | 1<<(_SIGSTOP-1))

//go:nosplit
//go:nowritebarrierrec
func redoxSigprocmask(how int32, set *sigset, oset *sigset) {
	old := redoxSigmaskLoad()
	if oset != nil {
		redoxSigsetSet(oset, old)
	}
	if set == nil {
		return
	}

	next := old
	mask := redoxSigsetGet(set)
	switch how {
	case _SIG_SETMASK:
		next = mask
	case _SIG_BLOCK:
		next = old | mask
	case _SIG_UNBLOCK:
		next = old &^ mask
	default:
		return
	}
	redoxSigmaskStore(next)
}

//go:nosplit
//go:nowritebarrierrec
func redoxSigmaskLoad() uint64 {
	ctl := redoxSigcontrolCurrent()
	w0 := atomic.Load64(&ctl.word[0])
	w1 := atomic.Load64(&ctl.word[1])
	allow := (w0 >> 32) | ((w1 >> 32) << 32)
	return ^allow
}

//go:nosplit
//go:nowritebarrierrec
func redoxSigmaskStore(mask uint64) {
	allow := ^mask | redoxSigmaskAllowAlways
	ctl := redoxSigcontrolCurrent()
	redoxSigmaskStoreWord(&ctl.word[0], (allow&0xffffffff)<<32)
	redoxSigmaskStoreWord(&ctl.word[1], allow&0xffffffff00000000)
}

//go:nosplit
//go:nowritebarrierrec
func redoxSigmaskStoreWord(word *uint64, allow uint64) {
	for {
		old := atomic.Load64(word)
		next := (old & 0xffffffff) | allow
		if atomic.Cas64(word, old, next) {
			return
		}
	}
}

//go:nosplit
func redoxSigsetGet(set *sigset) uint64 {
	return uint64(set.__bits[0]) | uint64(set.__bits[1])<<32
}

//go:nosplit
func redoxSigsetSet(set *sigset, mask uint64) {
	set.__bits[0] = uint32(mask)
	set.__bits[1] = uint32(mask >> 32)
}

func sysconf(name int32) int64 {
	return int64(sysvicall1(&libc_sysconf, uintptr(name)))
}

func usleep1(usec uint32)

//go:nosplit
func usleep_no_g(µs uint32) {
	usleep1(µs)
}

//go:nosplit
func usleep(µs uint32) {
	sysvicall1(&libc_usleep, uintptr(µs))
}

func walltime() (sec int64, nsec int32) {
	var ts timespec
	sysvicall2(&libc_clock_gettime, CLOCK_REALTIME, uintptr(unsafe.Pointer(&ts)))
	return ts.tv_sec, int32(ts.tv_nsec)
}

//go:nosplit
func write1(fd uintptr, buf unsafe.Pointer, nbyte int32) int32 {
	r1, err := sysvicall3Err(&libc_write, fd, uintptr(buf), uintptr(nbyte))
	KeepAlive(buf)
	if c := int32(r1); c >= 0 {
		return c
	}
	return -int32(err)
}

//go:nosplit
func pipe2(flags int32) (r, w int32, errno int32) {
	var p [2]int32
	_, e := sysvicall2Err(&libc_pipe2, uintptr(noescape(unsafe.Pointer(&p))), uintptr(flags))
	return p[0], p[1], int32(e)
}

//go:nosplit
func fcntl(fd, cmd, arg int32) (ret int32, errno int32) {
	r0, err := syscall_fcntl(uintptr(fd), uintptr(cmd), uintptr(arg))
	if err != 0 {
		errno = int32(err)
	} else {
		ret = int32(r0)
	}
	return
}
