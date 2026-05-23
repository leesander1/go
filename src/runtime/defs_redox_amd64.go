package runtime

const (
	_EINTR     = 0x4
	_EAGAIN    = 0xb
	_ETIMEDOUT = 0x6e

	_PROT_NONE  = 0x0
	_PROT_READ  = 0x4
	_PROT_WRITE = 0x2
	_PROT_EXEC  = 0x1

	_MAP_ANON    = 0x20
	_MAP_PRIVATE = 0x2
	_MAP_FIXED   = 0x4

	_MADV_DONTNEED   = 0x4
	_MADV_FREE       = 0x8
	_MADV_HUGEPAGE   = 0xe
	_MADV_NOHUGEPAGE = 0xf
	_MADV_COLLAPSE   = 0x19

	_SA_SIGINFO  = 0x2000000
	_SA_RESTART  = 0x8000000
	_SA_ONSTACK  = 0x4000000
	_SA_RESTORER = 0x4

	_SS_DISABLE  = 2
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
	_NSIG        = 32
	_SI_USER     = 0
	_UC_SIGMASK  = 0x01
	_UC_CPU      = 0x04

	_FPE_INTDIV = 0x1
	_FPE_INTOVF = 0x2
	_FPE_FLTDIV = 0x3
	_FPE_FLTOVF = 0x4
	_FPE_FLTUND = 0x5
	_FPE_FLTRES = 0x6
	_FPE_FLTINV = 0x7
	_FPE_FLTSUB = 0x8

	_SIGHUP    = 0x1
	_SIGINT    = 0x2
	_SIGQUIT   = 0x3
	_SIGILL    = 0x4
	_SIGTRAP   = 0x5
	_SIGABRT   = 0x6
	_SIGBUS    = 0x7
	_SIGFPE    = 0x8
	_SIGKILL   = 0x9
	_SIGUSR1   = 0xa
	_SIGSEGV   = 0xb
	_SIGUSR2   = 0xc
	_SIGPIPE   = 0xd
	_SIGALRM   = 0xe
	_SIGTERM   = 0xf
	_SIGSTKFLT = 0x10
	_SIGCHLD   = 0x11
	_SIGCONT   = 0x12
	_SIGSTOP   = 0x13
	_SIGTSTP   = 0x14
	_SIGTTIN   = 0x15
	_SIGTTOU   = 0x16
	_SIGURG    = 0x17
	_SIGXCPU   = 0x18
	_SIGXFSZ   = 0x19
	_SIGVTALRM = 0x1a
	_SIGPROF   = 0x1b
	_SIGWINCH  = 0x1c
	_SIGIO     = 0x1d
	_SIGPWR    = 0x1e
	_SIGSYS    = 0x1f

	_SIGRTMIN = 0x23

	_BUS_ADRALN = 0x1
	_BUS_ADRERR = 0x2
	_BUS_OBJERR = 0x3

	_SEGV_MAPERR = 0x1
	_SEGV_ACCERR = 0x2

	_ITIMER_REAL    = 0x0
	_ITIMER_VIRTUAL = 0x1
	_ITIMER_PROF    = 0x2

	_AF_UNIX    = 0x1
	_SOCK_DGRAM = 0x2

	_O_RDONLY   = 0x10000
	_O_WRONLY   = 0x20000
	_O_NONBLOCK = 0x40000

	_O_CLOEXEC = 0x1000000
	_O_CREAT   = 0x2000000
	_O_TRUNC   = 0x4000000
	_O_EXCL    = 0x8000000

	_O_DIRECTORY = 0x10000000
	_O_PATH      = 0x20000000
	_O_SYMLINK   = 0x40000000

	CLOCK_REALTIME  = 0
	CLOCK_MONOTONIC = 4

	PTHREAD_CREATE_DETACHED = 0

	_SC_NPROCESSORS_ONLN = 58
	_SC_PAGESIZE         = 30

	_MAXHOSTNAMELEN = 0x40
)

type pthread_t uintptr

type pthread_attr_t struct {
	__align [32]byte
}

type pthread_cond_t struct {
	__align [8]byte
}

type pthread_mutex_t struct {
	__align [12]byte
}

type sigset struct {
	__bits [2]uint32
}

type siginfo struct {
	si_signo  int32
	si_errno  int32
	si_code   int32
	si_pid    int32
	si_uid    int32
	si_addr   uintptr
	si_status int32
	si_value  uintptr
}

type sem_t struct {
	__align uintptr
}

type stackt struct {
	ss_sp    uintptr
	ss_flags int32
	ss_size  uintptr
}

type sigactiont struct {
	sa_handler  uintptr
	sa_flags    uint64
	sa_restorer uintptr
	sa_mask     uint64
}

type timespec struct {
	tv_sec  int64
	tv_nsec int64
}

//go:nosplit
func (ts *timespec) setNsec(ns int64) {
	ts.tv_sec = ns / 1e9
	ts.tv_nsec = ns % 1e9
}

type timeval struct {
	tv_sec  int64
	tv_usec int32
	__pad   [4]byte
}

func (tv *timeval) set_usec(x int32) {
	tv.tv_usec = x
}

type itimerval struct {
	it_interval timeval
	it_value    timeval
}

type mcontext struct {
	ymmupper [16][2]uint64
	fxsave   [29][2]uint64
	r15      uint64
	r14      uint64
	r13      uint64
	r12      uint64
	rbp      uint64
	rbx      uint64
	r11      uint64
	r10      uint64
	r9       uint64
	r8       uint64
	rax      uint64
	rcx      uint64
	rdx      uint64
	rsi      uint64
	rdi      uint64
	rflags   uint64
	rip      uint64
	rsp      uint64
}

type ucontext struct {
	__pad       [8]byte
	uc_link     *ucontext
	uc_stack    stackt
	uc_sigmask  sigset
	_sival      uintptr
	_sigcode    uint32
	_signum     uint32
	uc_mcontext mcontext
}

// this mirrors the ucontext struct
type sigcontext struct {
	pad_cgo_0  [8]byte
	uclink     *ucontext // or *sigcontext, depending on api
	ucstack    stackt
	ucsigmask  sigset
	sival      uintptr
	sigcode    uint32
	signum     uint32
	ucmcontext mcontext
}

type sockaddr_un struct {
	family uint16
	path   [108]byte
}

type mscratch struct {
	v [6]uintptr
}
