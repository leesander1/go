// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// _rt0_amd64_redox is the initial entry point from the operating system.
// Its job is to call relibc's initialization function, relibc_crt0.
TEXT _rt0_amd64_redox(SB),NOSPLIT,$0
    // The x86-64 System V ABI requires the stack to be 16-byte aligned
    // before a `call` instruction. The kernel provides an aligned stack,
    // but the `call` itself will push an 8-byte return address, misaligning it.
    // We adjust the stack pointer by 8 bytes to compensate.
    SUBQ $8, SP

    // Per the ABI, the first argument is passed in the RDI register.
    // The Rust function `relibc_crt0` expects the stack pointer (`sp: usize`).
    MOVQ SP, DI

    // Call relibc's startup function.
    CALL relibc_crt0(SB)

    // relibc_crt0 is a non-returning function (`-> !`).
    // If it ever returns, something is deeply wrong. Halt the processor.
    HLT

TEXT _rt0_amd64_redox_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)
