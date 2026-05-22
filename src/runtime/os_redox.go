// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/abi"
	"internal/runtime/atomic"
	"internal/goarch"
	"internal/runtime/sys"
	"unsafe"
)


// OS-specific state for a machine (m).
type mOS struct {
	waitsema uintptr // semaphore for parking on locks
	perrno   *int32  // pointer to tls errno
	// This is here to avoid using the G stack so the stack can move during the call.
	libcall libcall
	ts      timespec
	scratch mscratch
}

type libcFunc uintptr

//go:linkname asmsysvicall6x runtime.asmsysvicall6
var asmsysvicall6x libcFunc // name to take addr of asmsysvicall6

func asmsysvicall6() // declared for vet; do NOT call

//go:nosplit
func sysvicall0(fn *libcFunc) uintptr {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil // See comment in sys_darwin.go:libcCall
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 0
	libcall.args = uintptr(unsafe.Pointer(fn)) // it's unused but must be non-nil, otherwise crashes
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1
}

//go:nosplit
func sysvicall1(fn *libcFunc, a1 uintptr) uintptr {
	r1, _ := sysvicall1Err(fn, a1)
	return r1
}

// sysvicall1Err returns both the system call result and the errno value.
// This is used by sysvicall1, epoll_create1 and pipe.
//
//go:nosplit
func sysvicall1Err(fn *libcFunc, a1 uintptr) (r1, err uintptr) {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 1
	// TODO(rsc): Why is noescape necessary here and below?
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1, libcall.err
}

//go:nosplit
func sysvicall2(fn *libcFunc, a1, a2 uintptr) uintptr {
	r1, _ := sysvicall2Err(fn, a1, a2)
	return r1
}

//go:nosplit
//go:cgo_unsafe_args

// sysvicall2Err returns both the system call result and the errno value.
// This is used by sysvicall2 and pipe2.
func sysvicall2Err(fn *libcFunc, a1, a2 uintptr) (uintptr, uintptr) {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 2
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1, libcall.err
}

//go:nosplit
func sysvicall3(fn *libcFunc, a1, a2, a3 uintptr) uintptr {
	r1, _ := sysvicall3Err(fn, a1, a2, a3)
	return r1
}

//go:nosplit
//go:cgo_unsafe_args

// sysvicall3Err returns both the system call result and the errno value.
// This is used by sysvicall3 and write1.
func sysvicall3Err(fn *libcFunc, a1, a2, a3 uintptr) (r1, err uintptr) {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 3
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1, libcall.err
}

//go:nosplit
//go:cgo_unsafe_args
func sysvicall4(fn *libcFunc, a1, a2, a3, a4 uintptr) uintptr {
	r1, _ := sysvicall4Err(fn, a1, a2, a3, a4)
	return r1
}

//go:nosplit
//go:cgo_unsafe_args

// sysvicall4Err returns both the system call result and the errno value.
// This is used by sysvicall4 and epoll_wait.
func sysvicall4Err(fn *libcFunc, a1, a2, a3, a4 uintptr) (r1, err uintptr) {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 4
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1, libcall.err
}

//go:nosplit
//go:cgo_unsafe_args
func sysvicall5(fn *libcFunc, a1, a2, a3, a4, a5 uintptr) uintptr {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 5
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1
}

//go:nosplit
//go:cgo_unsafe_args
func sysvicall6(fn *libcFunc, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
	// Leave caller's PC/SP around for traceback.
	gp := getg()
	var mp *m
	if gp != nil {
		mp = gp.m
	}
	if mp != nil && mp.libcallsp == 0 {
		mp.libcallg.set(gp)
		mp.libcallpc = sys.GetCallerPC()
		// sp must be the last, because once async cpu profiler finds
		// all three values to be non-zero, it will use them
		mp.libcallsp = sys.GetCallerSP()
	} else {
		mp = nil
	}

	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(fn))
	libcall.n = 6
	libcall.args = uintptr(noescape(unsafe.Pointer(&a1)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	if mp != nil {
		mp.libcallsp = 0
	}
	return libcall.r1
}

var sigset_all = sigset{[2]uint32{^uint32(0), ^uint32(0)}}

//
// C library function declarations
//

func getCPUCount() int32 {
	n := int32(sysconf(_SC_NPROCESSORS_ONLN))
	if n < 1 {
		return 1
	}

	return n
}

func getPageSize() uintptr {
	n := int32(sysconf(_SC_PAGESIZE))
	if n <= 0 {
		return 0
	}
	return uintptr(n)
}

func osinit() {
	numCPUStartup = 1;

	physPageSize = 4096;
}

func tstart_sysvicall(newm *m) uint32

// May run with m.p==nil, so write barriers are not allowed.
//
//go:nowritebarrier
func newosproc(mp *m) {
	var (
		attr pthread_attr_t
		oset sigset
		tid  pthread_t
		ret  int32
		size uint64
	)

	if pthread_attr_init(&attr) != 0 {
		throw("pthread_attr_init")
	}
	// Allocate a new 2MB stack.
	if pthread_attr_setstack(&attr, 0, 0x200000) != 0 {
		throw("pthread_attr_setstack")
	}
	// Read back the allocated stack.
	if pthread_attr_getstack(&attr, unsafe.Pointer(&mp.g0.stack.hi), &size) != 0 {
		throw("pthread_attr_getstack")
	}
	mp.g0.stack.lo = mp.g0.stack.hi - uintptr(size)
	if pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED) != 0 {
		throw("pthread_attr_setdetachstate")
	}

	// Disable signals during create, so that the new thread starts
	// with signals disabled. It will enable them in minit.
	sigprocmask(_SIG_SETMASK, &sigset_all, &oset)
	ret = retryOnEAGAIN(func() int32 {
		return pthread_create(&tid, &attr, abi.FuncPCABI0(tstart_sysvicall), unsafe.Pointer(mp))
	})
	sigprocmask(_SIG_SETMASK, &oset, nil)
	if ret != 0 {
		print("runtime: failed to create new OS thread (have ", mcount(), " already; errno=", ret, ")\n")
		if ret == _EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("newosproc")
	}
}

func exitThread(wait *atomic.Uint32) {
	// We should never reach exitThread on Solaris because we let
	// libc clean up threads.
	throw("exitThread")
}

var urandom_dev = []byte("/scheme/rand\x00")

//go:nosplit
func readRandom(r []byte) int {
	// broken
	// print("openrandom\n")
	// var pp *byte
	// pp = &urandom_dev[0]
	// print("firstrandom ", uintptr(unsafe.Pointer(pp)), "\n")
	// fd := open(pp, _O_RDONLY, 0)
	// print("readrandom\n")
	// n := read(fd, unsafe.Pointer(&r[0]), int32(len(r)))
	// print("closerandom\n")
	// closefd(fd)
	// return int(n)
	return 0
}

func goenvs() {
    environs_uintptr := get_environ()
    envp := (**byte)(unsafe.Pointer(environs_uintptr))
    n := 0
    for *(**byte)(add(unsafe.Pointer(envp), uintptr(n)*goarch.PtrSize)) != nil {
        n++
    }
    envs = make([]string, n)
    for i := 0; i < n; i++ {
        envs[i] = gostring(argv_index(envp, int32(i)))
    }
}


// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024)
	mp.gsignal.m = mp
}

func miniterrno()

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, cannot allocate memory.
func minit() {
	asmcgocall(unsafe.Pointer(abi.FuncPCABI0(miniterrno)), unsafe.Pointer(&libc___errno))

	minitSignals()

	getg().m.procid = uint64(pthread_self())
}

// Called from dropm to undo the effect of an minit.
func unminit() {
	unminitSignals()
	getg().m.procid = 0
}

// Called from mexit, but not from dropm, to undo the effect of thread-owned
// resources in minit, semacreate, or elsewhere. Do not take locks after calling this.
//
// This always runs without a P, so //go:nowritebarrierrec is required.
//
//go:nowritebarrierrec
func mdestroy(mp *m) {
}

func sigtramp()

//go:nosplit
//go:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont

	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTART
	sa.sa_mask = 0                               // TODO
	if fn == abi.FuncPCABIInternal(sighandler) { // abi.FuncPCABIInternal(sighandler) matches the callers in signal_unix.go
		fn = abi.FuncPCABI0(sigtramp)
	}
	*((*uintptr)(unsafe.Pointer(&sa.sa_handler))) = fn
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func setsigstack(i uint32) {
	var sa sigactiont
	sigaction(i, nil, &sa)
	if sa.sa_flags&_SA_ONSTACK != 0 {
		return
	}
	sa.sa_flags |= _SA_ONSTACK
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func getsig(i uint32) uintptr {
	var sa sigactiont
	sigaction(i, nil, &sa)
	return *((*uintptr)(unsafe.Pointer(&sa.sa_handler)))
}

// setSignalstackSP sets the ss_sp field of a stackt.
//
//go:nosplit
func setSignalstackSP(s *stackt, sp uintptr) {
	*(*uintptr)(unsafe.Pointer(&s.ss_sp)) = sp
}

//go:nosplit
//go:nowritebarrierrec
func sigaddset(mask *sigset, i int) {
	mask.__bits[(i-1)/32] |= 1 << ((uint32(i) - 1) & 31)
}

func sigdelset(mask *sigset, i int) {
	mask.__bits[(i-1)/32] &^= 1 << ((uint32(i) - 1) & 31)
}

//go:nosplit
func (c *sigctxt) fixsigcode(sig uint32) {
}

func setProcessCPUProfiler(hz int32) {
	setProcessCPUProfilerTimer(hz)
}

func setThreadCPUProfiler(hz int32) {
	setThreadCPUProfilerHz(hz)
}

//go:nosplit
func validSIGPROF(mp *m, c *sigctxt) bool {
	return true
}

//go:nosplit
func semacreate(mp *m) {
	if mp.waitsema != 0 {
		return
	}

	var sem *sem_t

	// Call libc's malloc rather than malloc. This will
	// allocate space on the C heap. We can't call malloc
	// here because it could cause a deadlock.
	mp.libcall.fn = uintptr(unsafe.Pointer(&libc_malloc))
	mp.libcall.n = 1
	mp.scratch = mscratch{}
	mp.scratch.v[0] = unsafe.Sizeof(*sem)
	mp.libcall.args = uintptr(unsafe.Pointer(&mp.scratch))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&mp.libcall))
	sem = (*sem_t)(unsafe.Pointer(mp.libcall.r1))
	if sem_init(sem, 0, 0) != 0 {
		throw("sem_init")
	}
	mp.waitsema = uintptr(unsafe.Pointer(sem))
}

//go:nosplit
func semasleep(ns int64) int32 {
	mp := getg().m
	if ns >= 0 {
		// sem_timedwait requires an absolute timeout.
		var now timespec
		sysvicall2(&libc_clock_gettime, CLOCK_REALTIME, uintptr(unsafe.Pointer(&now)))
		mp.ts.tv_sec = now.tv_sec + ns/1e9
		mp.ts.tv_nsec = now.tv_nsec + ns%1e9

		// Handle nanosecond overflow.
		if mp.ts.tv_nsec >= 1e9 {
			mp.ts.tv_sec++
			mp.ts.tv_nsec -= 1e9
		}
		sysvicall2(&libc_sem_timedwait, mp.waitsema, uintptr(unsafe.Pointer(&mp.ts)))
		if *mp.perrno != 0 {
			if *mp.perrno == _ETIMEDOUT || *mp.perrno == _EAGAIN || *mp.perrno == _EINTR {
				return -1
			}
			throw("sem_timedwait")
		}
		return 0
	}
	for {
		if sysvicall1(&libc_sem_wait, mp.waitsema) == 0 {
			break
		}
		if *mp.perrno == _EINTR {
			continue
		}
		throw("sem_wait")
	}
	return 0
}

//go:nosplit
func semawakeup(mp *m) {
	if sem_post((*sem_t)(unsafe.Pointer(mp.waitsema))) != 0 {
		throw("sem_post")
	}
}


func osyield1()

//go:nosplit
func osyield_no_g() {
	osyield1()
}

//go:nosplit
func osyield() {
	sysvicall0(&libc_sched_yield)
}

//go:linkname executablePath os.executablePath
var executablePath string

// sigPerThreadSyscall is only used on linux, so we assign a bogus signal
// number.
const sigPerThreadSyscall = 1 << 31

//go:nosplit
func runPerThreadSyscall() {
	throw("runPerThreadSyscall only valid on linux")
}
