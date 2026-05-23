// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

//
// System calls for redox/amd64 are implemented in ../runtime/syscall_redox.go
//

#define SYS_CLOSE 0x20000006
#define SYS_DUP2 0x2010003f
#define SYS_FCNTL 0x20000037
#define SYS_WRITE 0x21000004

TEXT ·sysvicall6(SB),NOSPLIT,$0
	JMP	runtime·syscall_sysvicall6(SB)

TEXT ·rawSysvicall6(SB),NOSPLIT,$0
	JMP	runtime·syscall_rawsysvicall6(SB)

TEXT ·syscgocall6ABI0(SB),NOSPLIT,$0
	JMP	·syscgocall6(SB)

TEXT ·rawsyscgocall6ABI0(SB),NOSPLIT,$0
	JMP	·rawsyscgocall6(SB)

TEXT ·chdir(SB),NOSPLIT,$0
	JMP	runtime·syscall_chdir(SB)

TEXT ·chroot1(SB),NOSPLIT,$0
	JMP	runtime·syscall_chroot(SB)

TEXT ·closeFD(SB),NOSPLIT,$0-16
	MOVQ	$SYS_CLOSE, AX
	MOVQ	fd+0(FP), DI
	SYSCALL
	CMPQ	AX, $-4096
	JLS	closeok
	NEGQ	AX
	MOVQ	AX, err+8(FP)
	RET
closeok:
	MOVQ	$0, err+8(FP)
	RET

TEXT ·dup2child(SB),NOSPLIT,$0-32
	MOVQ	$SYS_DUP2, AX
	MOVQ	old+0(FP), DI
	MOVQ	new+8(FP), SI
	MOVQ	$0, DX
	MOVQ	$0, R10
	SYSCALL
	CMPQ	AX, $-4096
	JLS	dup2ok
	NEGQ	AX
	MOVQ	$0, val+16(FP)
	MOVQ	AX, err+24(FP)
	RET
dup2ok:
	MOVQ	AX, val+16(FP)
	MOVQ	$0, err+24(FP)
	RET

TEXT ·execve(SB),NOSPLIT,$0
	JMP	runtime·syscall_execve(SB)

TEXT ·execveStack(SB),NOSPLIT,$0
	JMP	runtime·syscall_execve_stack(SB)

TEXT ·exit(SB),NOSPLIT,$0
	JMP	runtime·syscall_exit(SB)

TEXT ·fcntl1(SB),NOSPLIT,$0-40
	MOVQ	$SYS_FCNTL, AX
	MOVQ	fd+0(FP), DI
	MOVQ	cmd+8(FP), SI
	MOVQ	arg+16(FP), DX
	SYSCALL
	CMPQ	AX, $-4096
	JLS	fcntlok
	NEGQ	AX
	MOVQ	$0, val+24(FP)
	MOVQ	AX, err+32(FP)
	RET
fcntlok:
	MOVQ	AX, val+24(FP)
	MOVQ	$0, err+32(FP)
	RET

TEXT ·forkx(SB),NOSPLIT,$0
	JMP	runtime·syscall_forkx(SB)

TEXT ·gethostname(SB),NOSPLIT,$0
	JMP	runtime·syscall_gethostname(SB)

TEXT ·getpid(SB),NOSPLIT,$0
	JMP	runtime·syscall_getpid(SB)

TEXT ·ioctl(SB),NOSPLIT,$0
	JMP	runtime·syscall_ioctl(SB)

TEXT ·RawSyscall(SB),NOSPLIT,$0
	JMP	runtime·syscall_rawsyscall(SB)

TEXT ·RawSyscall6(SB),NOSPLIT,$0
	JMP	runtime·syscall_rawsyscall6(SB)

TEXT ·setgid(SB),NOSPLIT,$0
	JMP	runtime·syscall_setgid(SB)

TEXT ·setgroups1(SB),NOSPLIT,$0
	JMP	runtime·syscall_setgroups(SB)

TEXT ·setrlimit1(SB),NOSPLIT,$0
	JMP	runtime·syscall_setrlimit(SB)

TEXT ·setsid(SB),NOSPLIT,$0
	JMP	runtime·syscall_setsid(SB)

TEXT ·setuid(SB),NOSPLIT,$0
	JMP	runtime·syscall_setuid(SB)

TEXT ·setpgid(SB),NOSPLIT,$0
	JMP	runtime·syscall_setpgid(SB)

TEXT ·Syscall(SB),NOSPLIT,$0
	JMP	runtime·syscall_syscall(SB)

TEXT ·wait4(SB),NOSPLIT,$0
	JMP	runtime·syscall_wait4(SB)

TEXT ·write1(SB),NOSPLIT,$0-40
	MOVQ	$SYS_WRITE, AX
	MOVQ	fd+0(FP), DI
	MOVQ	buf+8(FP), SI
	MOVQ	nbyte+16(FP), DX
	SYSCALL
	CMPQ	AX, $-4096
	JLS	writeok
	NEGQ	AX
	MOVQ	$0, n+24(FP)
	MOVQ	AX, err+32(FP)
	RET
writeok:
	MOVQ	AX, n+24(FP)
	MOVQ	$0, err+32(FP)
	RET
