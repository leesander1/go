// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix && !redox

package runtime

const (
	checkfdsDevNull  = "/dev/null"
	checkfdsOpenMode = 0x02 // O_RDWR
)
