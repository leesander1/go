// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package runtime

func checkfds() {
	if islibrary || isarchive {
		// If the program is actually a library, presumably being consumed by
		// another program, we don't want to mess around with the file
		// descriptors.
		return
	}

	const (
		// F_GETFD and EBADF are standard across all unixes we support, so we
		// define them here rather than in each of the OS specific files.
		F_GETFD = 0x01
		EBADF   = 0x09
	)

	devNull := []byte(checkfdsDevNull + "\x00")
	for i := 0; i < 3; i++ {
		ret, errno := fcntl(int32(i), F_GETFD, 0)
		if ret >= 0 {
			continue
		}

		if errno != EBADF {
			print("runtime: unexpected error while checking standard file descriptor ", i, ", errno=", errno, "\n")
			throw("cannot open standard fds")
		}

		if ret := open(&devNull[0], checkfdsOpenMode, 0); ret < 0 {
			print("runtime: standard file descriptor ", i, " closed, unable to open ", checkfdsDevNull, ", errno=", errno, "\n")
			throw("cannot open standard fds")
		} else if ret != int32(i) {
			print("runtime: opened unexpected file descriptor ", ret, " when attempting to open ", i, "\n")
			throw("cannot open standard fds")
		}
	}
}
