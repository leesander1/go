// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall

import (
	"sync"
	"unsafe"
)

const redoxExecStackSize = 1 << 20

var (
	redoxExecStackMu sync.Mutex
	redoxExecStack   [redoxExecStackSize]byte
)

func redoxExecStackAcquire() uintptr {
	redoxExecStackMu.Lock()
	return uintptr(unsafe.Pointer(&redoxExecStack[len(redoxExecStack)-1])) + 1
}

func redoxExecStackRelease() {
	redoxExecStackMu.Unlock()
}
