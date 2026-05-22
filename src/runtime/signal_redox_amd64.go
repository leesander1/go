// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/goarch"
	"unsafe"
)

type sigctxt struct {
	info *siginfo
	ctxt unsafe.Pointer
}

//go:nosplit
//go:nowritebarrierrec
func (c *sigctxt) regs() *sigcontext {
	// c.ctxt is a pointer to the full ucontext. since sigcontext is an
	// alias for that layout, we can cast it directly.
	return (*sigcontext)(c.ctxt)
}

// all register accessors must go through the ucmcontext field of the sigcontext.

func (c *sigctxt) rax() uint64 { return c.regs().ucmcontext.rax }
func (c *sigctxt) rbx() uint64 { return c.regs().ucmcontext.rbx }
func (c *sigctxt) rcx() uint64 { return c.regs().ucmcontext.rcx }
func (c *sigctxt) rdx() uint64 { return c.regs().ucmcontext.rdx }
func (c *sigctxt) rdi() uint64 { return c.regs().ucmcontext.rdi }
func (c *sigctxt) rsi() uint64 { return c.regs().ucmcontext.rsi }
func (c *sigctxt) rbp() uint64 { return c.regs().ucmcontext.rbp }
func (c *sigctxt) rsp() uint64 { return c.regs().ucmcontext.rsp }
func (c *sigctxt) r8() uint64  { return c.regs().ucmcontext.r8 }
func (c *sigctxt) r9() uint64  { return c.regs().ucmcontext.r9 }
func (c *sigctxt) r10() uint64 { return c.regs().ucmcontext.r10 }
func (c *sigctxt) r11() uint64 { return c.regs().ucmcontext.r11 }
func (c *sigctxt) r12() uint64 { return c.regs().ucmcontext.r12 }
func (c *sigctxt) r13() uint64 { return c.regs().ucmcontext.r13 }
func (c *sigctxt) r14() uint64 { return c.regs().ucmcontext.r14 }
func (c *sigctxt) r15() uint64 { return c.regs().ucmcontext.r15 }

//go:nosplit
//go:nowritebarrierrec
func (c *sigctxt) rip() uint64 { return c.regs().ucmcontext.rip }

func (c *sigctxt) rflags() uint64 { return c.regs().ucmcontext.rflags } // changed from eflags

// cs, fs, gs are not available in the relibc ucontext struct.
func (c *sigctxt) cs() uint64 { return uint64(0) }
func (c *sigctxt) fs() uint64 { return uint64(0) }
func (c *sigctxt) gs() uint64 { return uint64(0) }

func (c *sigctxt) sigcode() uint64 { return uint64(c.info.si_code) }
func (c *sigctxt) sigaddr() uint64 { return uint64(c.info.si_addr) }

func (c *sigctxt) set_rip(x uint64) { c.regs().ucmcontext.rip = x }
func (c *sigctxt) set_rsp(x uint64) { c.regs().ucmcontext.rsp = x }

func (c *sigctxt) set_sigcode(x uint64) { c.info.si_code = int32(x) }
func (c *sigctxt) set_sigaddr(x uint64) {
	*(*uintptr)(add(unsafe.Pointer(c.info), 2*goarch.PtrSize)) = uintptr(x)
}
