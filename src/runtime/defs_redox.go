// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

/*
Input to cgo.

GOARCH=amd64 go tool cgo -cdefs defs_redox.go >defs_redox_amd64.h
*/

package runtime

/*
#include <sys/epoll.h>
#include <sys/mman.h>
#include <sys/select.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <time.h>
#include <fcntl.h>
#include <semaphore.h>
#include <errno.h>
#include <signal.h>
#include <pthread.h>
#include <netdb.h>
#include <unistd.h>
*/
import "C"

const (
	_EINTR       = C.EINTR
	_EBADF       = C.EBADF
	_EFAULT      = C.EFAULT
	_EAGAIN      = C.EAGAIN
	_EBUSY       = C.EBUSY
	_ETIME       = C.ETIME
	_ETIMEDOUT   = C.ETIMEDOUT
	_EWOULDBLOCK = C.EWOULDBLOCK
	_EINPROGRESS = C.EINPROGRESS

	_PROT_NONE  = C.PROT_NONE
	_PROT_READ  = C.PROT_READ
	_PROT_WRITE = C.PROT_WRITE
	_PROT_EXEC  = C.PROT_EXEC

	_MAP_ANON    = C.MAP_ANON
	_MAP_PRIVATE = C.MAP_PRIVATE
	_MAP_FIXED   = C.MAP_FIXED

	_MADV_DONTNEED = C.MADV_DONTNEED

	_SA_SIGINFO = C.SA_SIGINFO
	_SA_RESTART = C.SA_RESTART
	_SA_ONSTACK = C.SA_ONSTACK

	_SS_DISABLE = C.SS_DISABLE

	// _FPE_INTDIV = C.FPE_INTDIV
	// _FPE_INTOVF = C.FPE_INTOVF
	// _FPE_FLTDIV = C.FPE_FLTDIV
	// _FPE_FLTOVF = C.FPE_FLTOVF
	// _FPE_FLTUND = C.FPE_FLTUND
	// _FPE_FLTRES = C.FPE_FLTRES
	// _FPE_FLTINV = C.FPE_FLTINV
	// _FPE_FLTSUB = C.FPE_FLTSUB

	_SIGHUP    = C.SIGHUP
	_SIGINT    = C.SIGINT
	_SIGQUIT   = C.SIGQUIT
	_SIGILL    = C.SIGILL
	_SIGTRAP   = C.SIGTRAP
	_SIGABRT   = C.SIGABRT
	_SIGFPE    = C.SIGFPE
	_SIGKILL   = C.SIGKILL
	_SIGBUS    = C.SIGBUS
	_SIGSEGV   = C.SIGSEGV
	_SIGSYS    = C.SIGSYS
	_SIGPIPE   = C.SIGPIPE
	_SIGALRM   = C.SIGALRM
	_SIGTERM   = C.SIGTERM
	_SIGURG    = C.SIGURG
	_SIGSTOP   = C.SIGSTOP
	_SIGTSTP   = C.SIGTSTP
	_SIGCONT   = C.SIGCONT
	_SIGCHLD   = C.SIGCHLD
	_SIGTTIN   = C.SIGTTIN
	_SIGTTOU   = C.SIGTTOU
	_SIGIO     = C.SIGIO
	_SIGXCPU   = C.SIGXCPU
	_SIGXFSZ   = C.SIGXFSZ
	_SIGVTALRM = C.SIGVTALRM
	_SIGPROF   = C.SIGPROF
	_SIGWINCH  = C.SIGWINCH
	_SIGUSR1   = C.SIGUSR1
	_SIGUSR2   = C.SIGUSR2

	// _BUS_ADRALN = C.BUS_ADRALN
	// _BUS_ADRERR = C.BUS_ADRERR
	// _BUS_OBJERR = C.BUS_OBJERR

	// _SEGV_MAPERR = C.SEGV_MAPERR
	// _SEGV_ACCERR = C.SEGV_ACCERR

	_ITIMER_REAL    = C.ITIMER_REAL
	_ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	_ITIMER_PROF    = C.ITIMER_PROF

	_O_RDONLY   = C.O_RDONLY
	_O_WRONLY   = C.O_WRONLY
	_O_NONBLOCK = C.O_NONBLOCK
	_O_CREAT    = C.O_CREAT
	_O_TRUNC    = C.O_TRUNC

	__SC_NPROCESSORS_ONLN = C._SC_NPROCESSORS_ONLN

	_PTHREAD_CREATE_DETACHED = C.PTHREAD_CREATE_DETACHED
)

type pthread_t C.pthread_t
type pthread_attr_t C.pthread_attr_t
type pthread_cond_t C.pthread_cond_t
type pthread_mutex_t C.pthread_mutex_t

type sigset C.sigset_t
type siginfo C.siginfo_t
type sem_t C.sem_t
type stackt C.stack_t
type sigactiont C.struct_sigaction

type timespec C.struct_timespec
type timeval C.struct_timeval
type itimerval C.struct_itimerval

type mcontext C.mcontext_t
type ucontext C.ucontext_t
