// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall

import "sync"

const (
	redoxExecStackGuard = 64 << 10
	redoxExecStackSize  = 16 << 20
)

var (
	redoxExecStackMu sync.Mutex
	redoxExecStack   uintptr
)

func redoxExecStackAcquire() (uintptr, Errno) {
	redoxExecStackMu.Lock()
	if redoxExecStack == 0 {
		stack, errno := redoxExecStackAlloc()
		if errno != 0 {
			redoxExecStackMu.Unlock()
			return 0, errno
		}
		redoxExecStack = stack
	}
	return redoxExecStack + redoxExecStackSize, 0
}

func redoxExecStackRelease() {
	redoxExecStackMu.Unlock()
}

func redoxExecStackAlloc() (uintptr, Errno) {
	total := uintptr(redoxExecStackGuard + redoxExecStackSize)
	base, err := mmap(0, total, PROT_NONE, MAP_ANON|MAP_PRIVATE, -1, 0)
	if err != nil {
		return 0, errnoFromError(err)
	}

	stackBase := base + redoxExecStackGuard
	stack, err := mmap(stackBase, redoxExecStackSize, PROT_READ|PROT_WRITE, MAP_ANON|MAP_PRIVATE|MAP_FIXED, -1, 0)
	if err != nil {
		munmap(base, total)
		return 0, errnoFromError(err)
	}
	if stack != stackBase {
		munmap(stack, redoxExecStackSize)
		munmap(base, total)
		return 0, EINVAL
	}
	return stackBase, 0
}

func errnoFromError(err error) Errno {
	if errno, ok := err.(Errno); ok {
		return errno
	}
	return EINVAL
}
