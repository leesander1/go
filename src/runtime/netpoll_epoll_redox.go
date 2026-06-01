// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/runtime/atomic"
	"unsafe"
)

var (
	epfd           int32         = -1 // epoll descriptor
	netpollEventFd uintptr            // timer fd for netpollBreak
	netpollWakeSig atomic.Uint32      // used to avoid duplicate calls of netpollBreak
)

var netpollWakeTimerPath = []byte("/scheme/time/4\x00")

//go:cgo_import_static _cgo_libc_epoll_create1
//go:cgo_import_static _cgo_libc_epoll_ctl
//go:cgo_import_static _cgo_libc_epoll_wait
//go:linkname libc_epoll_create1 _cgo_libc_epoll_create1
//go:linkname libc_epoll_ctl _cgo_libc_epoll_ctl
//go:linkname libc_epoll_wait _cgo_libc_epoll_wait

var (
	libc_epoll_create1,
	libc_epoll_ctl,
	libc_epoll_wait byte
)

type EpollEvent struct {
	Events    uint32
	Pad_cgo_0 [4]byte
	Data      [8]byte // unaligned uintptr
	X_pad     uint64
}

const (
	AT_FDCWD = -0x64

	ENOENT = 0x2
	ENOSYS = 0x26

	EPOLLIN       = 0x1
	EPOLLOUT      = 0x4
	EPOLLERR      = 0x8
	EPOLLHUP      = 0x10
	EPOLLRDHUP    = 0x2000
	EPOLLET       = 0x80000000
	EPOLL_CLOEXEC = 0x01000000
	EPOLL_CTL_ADD = 0x1
	EPOLL_CTL_DEL = 0x2
	EPOLL_CTL_MOD = 0x3
)

func epoll_create1(flags int32) (r1 int32, err int32) {
	ret, errno := cgocaller1(unsafe.Pointer(&libc_epoll_create1), uintptr(flags))
	if errno != 0 {
		err = errno
	} else {
		r1 = int32(ret)
	}
	return
}

func epoll_ctl(epfd int32, op int32, fd int32, event *EpollEvent) int32 {
	if _, errno := cgocaller4(unsafe.Pointer(&libc_epoll_ctl), uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event))); errno != 0 {
		return errno
	}
	return 0
}

func epoll_wait(epfd int32, events *EpollEvent, maxevents int32, timeout int32) (r1 int32, err int32) {
	ret, errno := cgocaller4(unsafe.Pointer(&libc_epoll_wait), uintptr(epfd), uintptr(unsafe.Pointer(events)), uintptr(maxevents), uintptr(timeout))
	if errno != 0 {
		err = errno
	} else {
		r1 = int32(ret)
	}
	return
}

func redoxNormalizeEpollEvent(ev *EpollEvent) bool {
	data := *(*uintptr)(unsafe.Pointer(&ev.Data))
	if data >= uintptr(1)<<tagBits {
		return true
	}
	if ev.X_pad == 0 {
		return false
	}
	var events uint32
	if data&1 != 0 {
		events |= EPOLLIN
	}
	if data&2 != 0 {
		events |= EPOLLOUT
	}
	ev.Events = events
	*(*uintptr)(unsafe.Pointer(&ev.Data)) = uintptr(ev.X_pad)
	ev.X_pad = 0
	return true
}

func netpollinit() {
	var errno int32
	epfd, errno = epoll_create1(EPOLL_CLOEXEC)
	if errno != 0 {
		println("runtime: epollcreate failed with", errno)
		throw("runtime: netpollinit failed")
	}
	timerfd := open(&netpollWakeTimerPath[0], _O_RDONLY|_O_WRONLY|_O_CLOEXEC, 0)
	if timerfd < 0 {
		println("runtime: netpoll timer open failed")
		syscall_close(epfd)
		epfd = -1
		return
	}
	ev := EpollEvent{
		Events: EPOLLIN,
	}
	netpollEventFd = uintptr(timerfd)
	*(**uintptr)(unsafe.Pointer(&ev.Data)) = &netpollEventFd
	errno = epoll_ctl(epfd, EPOLL_CTL_ADD, timerfd, &ev)
	if errno != 0 {
		println("runtime: epollctl failed with", errno)
		syscall_close(timerfd)
		syscall_close(epfd)
		netpollEventFd = 0
		epfd = -1
	}
}

func netpollIsPollDescriptor(fd uintptr) bool {
	return fd == uintptr(epfd) || fd == netpollEventFd
}

func netpollopen(fd uintptr, pd *pollDesc) int32 {
	if epfd == -1 {
		return ENOSYS
	}
	var ev EpollEvent
	ev.Events = EPOLLIN | EPOLLOUT | EPOLLRDHUP | EPOLLET
	tp := taggedPointerPack(unsafe.Pointer(pd), pd.fdseq.Load())
	*(*taggedPointer)(unsafe.Pointer(&ev.Data)) = tp
	return epoll_ctl(epfd, EPOLL_CTL_ADD, int32(fd), &ev)
}

func netpollclose(fd uintptr) int32 {
	if epfd == -1 {
		return 0
	}
	var ev EpollEvent
	return epoll_ctl(epfd, EPOLL_CTL_DEL, int32(fd), &ev)
}

func netpollarm(pd *pollDesc, mode int) {
	throw("runtime: unused")
}

// netpollBreak interrupts an epollwait.
func netpollBreak() {
	if epfd == -1 || netpollEventFd == 0 {
		return
	}

	// Failing to cas indicates there is an in-flight wakeup, so we're done here.
	if !netpollWakeSig.CompareAndSwap(0, 1) {
		return
	}

	var ts timespec
	sysvicall2(&libc_clock_gettime, CLOCK_MONOTONIC, uintptr(unsafe.Pointer(&ts)))
	tsSize := int32(unsafe.Sizeof(ts))
	for {
		n := write(netpollEventFd, noescape(unsafe.Pointer(&ts)), tsSize)
		if n == tsSize {
			break
		}
		if n == -_EINTR {
			continue
		}
		if n == -_EAGAIN {
			return
		}
		println("runtime: netpollBreak write failed with", -n)
		throw("runtime: netpollBreak write failed")
	}
}

// netpoll checks for ready network connections.
// Returns a list of goroutines that become runnable,
// and a delta to add to netpollWaiters.
// This must never return an empty list with a non-zero delta.
//
// delay < 0: blocks indefinitely
// delay == 0: does not block, just polls
// delay > 0: block for up to that many nanoseconds
func netpoll(delay int64) (gList, int32) {
	if epfd == -1 {
		return gList{}, 0
	}
	var waitms int32
	if delay < 0 {
		waitms = -1
	} else if delay == 0 {
		waitms = 0
	} else if delay < 1e6 {
		waitms = 1
	} else if delay < 1e15 {
		waitms = int32(delay / 1e6)
		if waitms > 10 {
			// Redox can miss a newly earlier Go timer while an M is parked in the
			// event queue. Keep timed polls short so the scheduler rechecks all
			// timer heaps promptly even when the poll deadline is stale.
			waitms = 10
		}
	} else {
		// An arbitrary cap on how long to wait for a timer.
		// 1e9 ms == ~11.5 days.
		waitms = 1e9
	}
	var events [128]EpollEvent
retry:
	n, errno := epoll_wait(epfd, &events[0], int32(len(events)), waitms)
	if errno != 0 {
		if errno != _EINTR {
			println("runtime: epollwait on fd", epfd, "failed with", errno)
			throw("runtime: netpoll failed")
		}
		// If a timed sleep was interrupted, just return to
		// recalculate how long we should sleep now.
		if waitms > 0 {
			return gList{}, 0
		}
		goto retry
	}
	var toRun gList
	delta := int32(0)
	for i := int32(0); i < n; i++ {
		ev := events[i]
		if !redoxNormalizeEpollEvent(&ev) {
			continue
		}
		if ev.Events == 0 {
			continue
		}

		if *(**uintptr)(unsafe.Pointer(&ev.Data)) == &netpollEventFd {
			if ev.Events != EPOLLIN {
				println("runtime: netpoll: eventfd ready for", ev.Events)
				throw("runtime: netpoll: eventfd ready for something unexpected")
			}
			netpollWakeSig.Store(0)
			continue
		}

		var mode int32
		if ev.Events&(EPOLLIN|EPOLLRDHUP|EPOLLHUP|EPOLLERR) != 0 {
			mode += 'r'
		}
		if ev.Events&(EPOLLOUT|EPOLLHUP|EPOLLERR) != 0 {
			mode += 'w'
		}
		if mode != 0 {
			tp := *(*taggedPointer)(unsafe.Pointer(&ev.Data))
			pd := (*pollDesc)(tp.pointer())
			if pd == nil {
				continue
			}
			tag := tp.tag()
			if pd.fdseq.Load() == tag {
				pd.setEventErr(ev.Events == EPOLLERR, tag)
				delta += netpollready(&toRun, pd, mode)
			}
		}
	}
	return toRun, delta
}
