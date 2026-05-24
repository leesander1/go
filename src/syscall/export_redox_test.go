// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package syscall

import "unsafe"

func Tcgetpgrp(fd int) (pgid int32, err error) {
	if err := ioctlPtr(uintptr(fd), uintptr(TIOCGPGRP), unsafe.Pointer(&pgid)); err != 0 {
		return -1, err
	}
	return pgid, nil
}

func Tcsetpgrp(fd int, pgid int32) (err error) {
	if err := ioctlPtr(uintptr(fd), uintptr(TIOCSPGRP), unsafe.Pointer(&pgid)); err != 0 {
		return err
	}
	return nil
}
