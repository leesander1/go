// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || solaris

package syscall

func redoxExecStackAcquire() uintptr {
	return 0
}

func redoxExecStackRelease() {
}
