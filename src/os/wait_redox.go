// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"runtime"
	"syscall"
)

// blockUntilWaitable observes the child's exit without releasing its PID.
// pidWait can then exclude concurrent signals before consuming the status.
func (p *Process) blockUntilWaitable() (bool, error) {
	var status syscall.WaitStatus
	_, err := ignoringEINTR2(func() (int, error) {
		return syscall.Wait4(p.Pid, &status, syscall.WNOWAIT, nil)
	})
	runtime.KeepAlive(p)
	if err != nil {
		return false, NewSyscallError("wait", err)
	}
	return true, nil
}
