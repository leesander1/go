// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc

#include "textflag.h"

//
// System calls for amd64, Redox are implemented in runtime/syscall_redox.go
//

TEXT ·sysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·syscgocall6ABI0(SB)

TEXT ·rawSysvicall6(SB),NOSPLIT,$0-88
	JMP	syscall·rawsyscgocall6ABI0(SB)
