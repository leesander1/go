// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// System calls for Redox are implemented in runtime/syscall_redox.go

TEXT ·syscall6(SB),NOSPLIT,$0-88
	JMP	syscall·syscgocall6ABI0(SB)

TEXT ·rawSyscall6(SB),NOSPLIT,$0-88
	JMP	syscall·rawsyscgocall6ABI0(SB)
