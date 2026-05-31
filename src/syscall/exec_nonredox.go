// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package syscall

func redoxExecStackAcquire() (uintptr, Errno) {
	return 0, 0
}

func redoxExecStackRelease() {
}
