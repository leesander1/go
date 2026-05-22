// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

import (
	"sync/atomic"
)

// origRlimitNofile, if non-nil, is the original soft RLIMIT_NOFILE.
var origRlimitNofile atomic.Pointer[Rlimit]

func Setrlimit(resource int, rlim *Rlimit) error {
	return nil
}
