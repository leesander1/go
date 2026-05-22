// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/runtime/atomic"
	"unsafe"
)

var (
	epfd              int32         = -1 // epoll descriptor
	netpollEventFd    uintptr            // eventfd for netpollBreak
	netpollWakeWriter uintptr            // the WRITE end of the pipe
	netpollWakeSig    atomic.Uint32      // used to avoid duplicate calls of netpollBreak
)

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
	Events uint32
	Data   [8]byte // unaligned uintptr
	_Pad   [8]byte
}

const (
	AT_FDCWD = -0x64

	ENOENT = 0x2

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
	print("bout to run call libc_epoll_create1\n")
	ret, errno := cgocaller1(unsafe.Pointer(&libc_epoll_create1), uintptr(flags));
	if errno != 0 {
		err = errno
	} else {
		r1 = int32(ret)
		print("libc_epoll_create1: epfd ", r1, "\n")
	}
	return
}

func epoll_ctl(epfd int32, op int32, fd int32, event *EpollEvent) int32 {
	print("bout to run call libc_epoll_ctl: epfd ", epfd, " op ", op, " fd ", fd, "\n")
	if _, errno := cgocaller4(unsafe.Pointer(&libc_epoll_ctl), uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event))); errno != 0 {
		return errno
	}
	return 0
}

func epoll_wait(epfd int32, events *EpollEvent, maxevents int32, timeout int32) (r1 int32, err int32) {
	print("bout to run call libc_epoll_wait\n")
	ret, errno := cgocaller4(unsafe.Pointer(&libc_epoll_wait), uintptr(epfd), uintptr(unsafe.Pointer(events)), uintptr(maxevents), uintptr(timeout));
	if errno != 0 {
		err = errno
	} else {
		r1 = int32(ret)
	}
	return
}

func netpollinit() {
	var errno int32
	epfd, errno = epoll_create1(EPOLL_CLOEXEC)
	if errno != 0 {
		println("runtime: epollcreate failed with", errno)
		throw("runtime: netpollinit failed")
	}
	r, w, errno := nonblockingPipe()
	if errno != 0 {
		println("runtime: pipe failed with", errno)
		throw("runtime: netpollinit failed")
	}
	ev := EpollEvent{
		Events: EPOLLIN,
	}
	netpollEventFd = uintptr(r)
	netpollWakeWriter = uintptr(w)
	*(**uintptr)(unsafe.Pointer(&ev.Data)) = &netpollEventFd
	errno = epoll_ctl(epfd, EPOLL_CTL_ADD, r, &ev)
	if errno != 0 {
		println("runtime: epollctl failed with", errno)
		// throw("runtime: epollctl failed")
	}
}

func netpollIsPollDescriptor(fd uintptr) bool {
	return fd == uintptr(epfd) || fd == netpollEventFd
}

func netpollopen(fd uintptr, pd *pollDesc) int32 {
	var ev EpollEvent
	ev.Events = EPOLLIN | EPOLLOUT | EPOLLRDHUP | EPOLLET
	tp := taggedPointerPack(unsafe.Pointer(pd), pd.fdseq.Load())
	*(*taggedPointer)(unsafe.Pointer(&ev.Data)) = tp
	return epoll_ctl(epfd, EPOLL_CTL_ADD, int32(fd), &ev)
}

func netpollclose(fd uintptr) int32 {
	var ev EpollEvent
	return epoll_ctl(epfd, EPOLL_CTL_DEL, int32(fd), &ev)
}

func netpollarm(pd *pollDesc, mode int) {
	throw("runtime: unused")
}

// netpollBreak interrupts an epollwait.
func netpollBreak() {
	// Failing to cas indicates there is an in-flight wakeup, so we're done here.
	if !netpollWakeSig.CompareAndSwap(0, 1) {
		return
	}

	var one uint64 = 1
	oneSize := int32(unsafe.Sizeof(one))
	for {
		n := write(netpollWakeWriter, noescape(unsafe.Pointer(&one)), oneSize)
		if n == oneSize {
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
		if ev.Events == 0 {
			continue
		}

		if *(**uintptr)(unsafe.Pointer(&ev.Data)) == &netpollEventFd {
			if ev.Events != EPOLLIN {
				println("runtime: netpoll: eventfd ready for", ev.Events)
				throw("runtime: netpoll: eventfd ready for something unexpected")
			}
			if delay != 0 {
				// netpollBreak could be picked up by a
				// nonblocking poll. Only read the 8-byte
				// integer if blocking.
				// Since EFD_SEMAPHORE was not specified,
				// the eventfd counter will be reset to 0.
				var one uint64
				read(int32(netpollEventFd), noescape(unsafe.Pointer(&one)), int32(unsafe.Sizeof(one)))
				netpollWakeSig.Store(0)
			}
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
			tag := tp.tag()
			if pd.fdseq.Load() == tag {
				pd.setEventErr(ev.Events == EPOLLERR, tag)
				delta += netpollready(&toRun, pd, mode)
			}
		}
	}
	return toRun, delta
}
