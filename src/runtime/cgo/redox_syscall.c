// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

#include <fcntl.h>
#include <grp.h>
#include <dirent.h>
#include <signal.h>
#include <errno.h>
#include <time.h>
#include <unistd.h>
#include <pthread.h>
#include <semaphore.h>
#include <sched.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/resource.h>
#include <sys/utsname.h>
#include <sys/uio.h>
#include <sys/mman.h>
#include <sys/epoll.h>
#include <sys/wait.h>
#include <sys/ioctl.h>

#include "libcgo.h"

// argset_t matches runtime/cgocall.go:argset.
typedef struct {
	uintptr_t* args;
	uintptr_t retval;
	int error;
} argset_t;

extern uintptr_t redox_getdents_v0(uintptr_t fd, void* buf, uintptr_t count, uint64_t opaque);

// libc backed posix-compliant syscalls.

// SET_RETVAL is a macro to handle the standard POSIX return convention.
// Most syscalls return -1 on error and set the global `errno`.
// This macro checks for that condition and sets the return value accordingly.
#define SET_RETVAL(fn) \
  uintptr_t ret = (uintptr_t) fn ; \
  if (ret == (uintptr_t) -1) {	   \
    x->error =              errno; \
  } else                           \
    x->retval = ret

// SET_PTHREAD_RETVAL handles the return convention for most pthread functions.
// They return 0 on success and an error number on failure.
#define SET_PTHREAD_RETVAL(fn) \
	x->retval = (uintptr_t) fn;


// --- File Descriptor Operations ---

void
_cgo_libc_open(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	int flags = (int)x->args[1];
	mode_t mode = (mode_t)x->args[2];
	SET_RETVAL(open(pathname, flags, mode));
}

void
_cgo_libc_openat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	int flags = (int)x->args[2];
	mode_t mode = (mode_t)x->args[3];
	SET_RETVAL(openat(dirfd, pathname, flags, mode));
}

void
_cgo_libc_close(argset_t* x) {
	int fd = (int)x->args[0];
	SET_RETVAL(close(fd));
}

void
_cgo_libc_write(argset_t* x) {
	int fd = (int)x->args[0];
	const void* buf = (const void*)x->args[1];
	size_t count = (size_t)x->args[2];
	SET_RETVAL(write(fd, buf, count));
}

void
_cgo_libc_read(argset_t* x) {
	int fd = (int)x->args[0];
	void* buf = (void*)x->args[1];
	size_t count = (size_t)x->args[2];
	SET_RETVAL(read(fd, buf, count));
}

void
_cgo_libc_pread(argset_t* x) {
	int fd = (int)x->args[0];
	void* buf = (void*)x->args[1];
	size_t count = (size_t)x->args[2];
	off_t offset = (off_t)x->args[3];
	SET_RETVAL(pread(fd, buf, count, offset));
}

void
_cgo_libc_pwrite(argset_t* x) {
	int fd = (int)x->args[0];
	const void* buf = (const void*)x->args[1];
	size_t count = (size_t)x->args[2];
	off_t offset = (off_t)x->args[3];
	SET_RETVAL(pwrite(fd, buf, count, offset));
}

void
_cgo_libc_pipe2(argset_t* x) {
	int* p = (int*)x->args[0];
	int flags = (int)x->args[1];
	SET_RETVAL(pipe2(p, flags));
}

void
_cgo_libc_lseek(argset_t* x) {
	int fd = (int)x->args[0];
	off_t offset = (off_t)x->args[1];
	int whence = (int)x->args[2];
	SET_RETVAL(lseek(fd, offset, whence));
}

void
_cgo_libc_dup(argset_t* x) {
	int oldfd = (int)x->args[0];
	SET_RETVAL(dup(oldfd));
}

void
_cgo_libc_dup2(argset_t* x) {
	int oldfd = (int)x->args[0];
	int newfd = (int)x->args[1];
	SET_RETVAL(dup2(oldfd, newfd));
}

void
_cgo_libc_fcntl(argset_t* x) {
	int fd = (int)x->args[0];
	int cmd = (int)x->args[1];
	uintptr_t arg = (uintptr_t)x->args[2]; // Can be int, struct*, etc.
	SET_RETVAL(fcntl(fd, cmd, arg));
}

void
_cgo_libc_fsync(argset_t* x) {
	int fd = (int)x->args[0];
	SET_RETVAL(fsync(fd));
}

void
_cgo_libc_ftruncate(argset_t* x) {
	int fd = (int)x->args[0];
	off_t length = (off_t)x->args[1];
	SET_RETVAL(ftruncate(fd, length));
}

void
_cgo_libc_posix_getdents(argset_t* x) {
	int fd = (int)x->args[0];
	void* buf = (void*)x->args[1];
	size_t count = (size_t)x->args[2];
	int flag = (int)x->args[3];
	SET_RETVAL(posix_getdents(fd, buf, count, flag));
}

void
_cgo_libc_redox_getdents(argset_t* x) {
	uintptr_t fd = x->args[0];
	void* buf = (void*)x->args[1];
	uintptr_t count = x->args[2];
	uint64_t opaque = (uint64_t)x->args[3];
	x->retval = redox_getdents_v0(fd, buf, count, opaque);
}

void
_cgo_libc_writev(argset_t* x) {
    int fd = (int)x->args[0];
    const struct iovec* iov = (const struct iovec*)x->args[1];
    int iovcnt = (int)x->args[2];
    SET_RETVAL(writev(fd, iov, iovcnt));
}

void
_cgo_libc_sendfile(argset_t* x) {
	x->retval = (uintptr_t)-1;
	x->error = ENOSYS;
}

void
_cgo_libc_select(argset_t* x) {
	int nfds = (int)x->args[0];
	fd_set* readfds = (fd_set*)x->args[1];
	fd_set* writefds = (fd_set*)x->args[2];
	fd_set* exceptfds = (fd_set*)x->args[3];
	struct timeval* timeout = (struct timeval*)x->args[4];
	SET_RETVAL(select(nfds, readfds, writefds, exceptfds, timeout));
}

void
_cgo_libc_ioctl(argset_t* x) {
	int fd = (int)x->args[0];
	unsigned long request = (unsigned long)x->args[1];
	uintptr_t argp = (uintptr_t)x->args[2]; // Variadic argument
	SET_RETVAL(ioctl(fd, request, (void*)argp));
}

// --- Filesystem Operations ---

void
_cgo_libc_stat(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	struct stat* statbuf = (struct stat*)x->args[1];
	SET_RETVAL(stat(pathname, statbuf));
}

void
_cgo_libc_fstat(argset_t* x) {
	int fd = (int)x->args[0];
	struct stat* statbuf = (struct stat*)x->args[1];
	SET_RETVAL(fstat(fd, statbuf));
}

void
_cgo_libc_fstatat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	struct stat* statbuf = (struct stat*)x->args[2];
	int flags = (int)x->args[3];
	SET_RETVAL(fstatat(dirfd, pathname, statbuf, flags));
}

void
_cgo_libc_lstat(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	struct stat* statbuf = (struct stat*)x->args[1];
	SET_RETVAL(lstat(pathname, statbuf));
}

void
_cgo_libc_access(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	int mode = (int)x->args[1];
	SET_RETVAL(access(pathname, mode));
}

void
_cgo_libc_chdir(argset_t* x) {
	const char* path = (const char*)x->args[0];
	SET_RETVAL(chdir(path));
}

void
_cgo_libc_fchdir(argset_t* x) {
	int fd = (int)x->args[0];
	SET_RETVAL(fchdir(fd));
}

void
_cgo_libc_getcwd(argset_t* x) {
	char* buf = (char*)x->args[0];
	size_t size = (size_t)x->args[1];
	SET_RETVAL(getcwd(buf, size));
}

void
_cgo_libc_chmod(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	mode_t mode = (mode_t)x->args[1];
	SET_RETVAL(chmod(pathname, mode));
}

void
_cgo_libc_fchmodat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	mode_t mode = (mode_t)x->args[2];
	int flags = (int)x->args[3];
	SET_RETVAL(fchmodat(dirfd, pathname, mode, flags));
}

void
_cgo_libc_fchmod(argset_t* x) {
	int fd = (int)x->args[0];
	mode_t mode = (mode_t)x->args[1];
	SET_RETVAL(fchmod(fd, mode));
}

void
_cgo_libc_chown(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	uid_t owner = (uid_t)x->args[1];
	gid_t group = (gid_t)x->args[2];
	SET_RETVAL(chown(pathname, owner, group));
}

void
_cgo_libc_fchownat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	uid_t owner = (uid_t)x->args[2];
	gid_t group = (gid_t)x->args[3];
	int flags = (int)x->args[4];
	SET_RETVAL(fchownat(dirfd, pathname, owner, group, flags));
}

void
_cgo_libc_fchown(argset_t* x) {
	int fd = (int)x->args[0];
	uid_t owner = (uid_t)x->args[1];
	gid_t group = (gid_t)x->args[2];
	SET_RETVAL(fchown(fd, owner, group));
}

void
_cgo_libc_lchown(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	uid_t owner = (uid_t)x->args[1];
	gid_t group = (gid_t)x->args[2];
	SET_RETVAL(lchown(pathname, owner, group));
}

void
_cgo_libc_chroot(argset_t* x) {
	const char* path = (const char*)x->args[0];
	SET_RETVAL(chroot(path));
}

void
_cgo_libc_link(argset_t* x) {
	const char* oldpath = (const char*)x->args[0];
	const char* newpath = (const char*)x->args[1];
	SET_RETVAL(link(oldpath, newpath));
}

void
_cgo_libc_linkat(argset_t* x) {
	int olddirfd = (int)x->args[0];
	const char* oldpath = (const char*)x->args[1];
	int newdirfd = (int)x->args[2];
	const char* newpath = (const char*)x->args[3];
	int flags = (int)x->args[4];
	SET_RETVAL(linkat(olddirfd, oldpath, newdirfd, newpath, flags));
}

void
_cgo_libc_readlink(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	char* buf = (char*)x->args[1];
	size_t bufsiz = (size_t)x->args[2];
	SET_RETVAL(readlink(pathname, buf, bufsiz));
}

void
_cgo_libc_readlinkat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	char* buf = (char*)x->args[2];
	size_t bufsiz = (size_t)x->args[3];
	SET_RETVAL(readlinkat(dirfd, pathname, buf, bufsiz));
}

void
_cgo_libc_rename(argset_t* x) {
	const char* oldpath = (const char*)x->args[0];
	const char* newpath = (const char*)x->args[1];
	SET_RETVAL(rename(oldpath, newpath));
}

void
_cgo_libc_renameat(argset_t* x) {
	int olddirfd = (int)x->args[0];
	const char* oldpath = (const char*)x->args[1];
	int newdirfd = (int)x->args[2];
	const char* newpath = (const char*)x->args[3];
	SET_RETVAL(renameat(olddirfd, oldpath, newdirfd, newpath));
}

void
_cgo_libc_rmdir(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	SET_RETVAL(rmdir(pathname));
}

void
_cgo_libc_symlink(argset_t* x) {
	const char* target = (const char*)x->args[0];
	const char* linkpath = (const char*)x->args[1];
	SET_RETVAL(symlink(target, linkpath));
}

void
_cgo_libc_symlinkat(argset_t* x) {
	const char* target = (const char*)x->args[0];
	int newdirfd = (int)x->args[1];
	const char* linkpath = (const char*)x->args[2];
	SET_RETVAL(symlinkat(target, newdirfd, linkpath));
}

void
_cgo_libc_unlink(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	SET_RETVAL(unlink(pathname));
}

void
_cgo_libc_unlinkat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	int flags = (int)x->args[2];
	SET_RETVAL(unlinkat(dirfd, pathname, flags));
}

void
_cgo_libc_mkdir(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	mode_t mode = (mode_t)x->args[1];
	SET_RETVAL(mkdir(pathname, mode));
}

void
_cgo_libc_mkdirat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	mode_t mode = (mode_t)x->args[2];
	SET_RETVAL(mkdirat(dirfd, pathname, mode));
}

void
_cgo_libc_mknod(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	mode_t mode = (mode_t)x->args[1];
	dev_t dev = (dev_t)x->args[2];
	SET_RETVAL(mknod(pathname, mode, dev));
}

void
_cgo_libc_pathconf(argset_t* x) {
	const char* path = (const char*)x->args[0];
	int name = (int)x->args[1];
	SET_RETVAL(pathconf(path, name));
}

void
_cgo_libc_fpathconf(argset_t* x) {
	int fd = (int)x->args[0];
	int name = (int)x->args[1];
	SET_RETVAL(fpathconf(fd, name));
}

void
_cgo_libc_truncate(argset_t* x) {
	const char* path = (const char*)x->args[0];
	off_t length = (off_t)x->args[1];
	SET_RETVAL(truncate(path, length));
}

void
_cgo_libc_umask(argset_t* x) {
	mode_t mask = (mode_t)x->args[0];
	// umask() is always successful. The return value is the old mask.
	x->retval = (uintptr_t)umask(mask);
}

void
_cgo_libc_utimes(argset_t* x) {
	const char* filename = (const char*)x->args[0];
	const struct timeval* times = (const struct timeval*)x->args[1];
	SET_RETVAL(utimes(filename, times));
}

void
_cgo_libc_utimensat(argset_t* x) {
	int dirfd = (int)x->args[0];
	const char* pathname = (const char*)x->args[1];
	const struct timespec* times = (const struct timespec*)x->args[2];
	int flags = (int)x->args[3];
	SET_RETVAL(utimensat(dirfd, pathname, times, flags));
}

void
_cgo_libc_sync(argset_t* x) {
	sync();
	x->retval = 0; // sync() is always successful
}

// --- Process Information ---

void
_cgo_libc_getpid(argset_t* x) {
	x->retval = (uintptr_t)getpid();
}

void
_cgo_libc_getppid(argset_t* x) {
	x->retval = (uintptr_t)getppid();
}

void
_cgo_libc_getuid(argset_t* x) {
	x->retval = (uintptr_t)getuid();
}

void
_cgo_libc_geteuid(argset_t* x) {
	x->retval = (uintptr_t)geteuid();
}

void
_cgo_libc_getgid(argset_t* x) {
	x->retval = (uintptr_t)getgid();
}

void
_cgo_libc_getegid(argset_t* x) {
	x->retval = (uintptr_t)getegid();
}

void
_cgo_libc_getgroups(argset_t* x) {
	int size = (int)x->args[0];
	gid_t* list = (gid_t*)x->args[1];
	SET_RETVAL(getgroups(size, list));
}

// --- Process Management ---

void
_cgo_libc_fork(argset_t* x) {
	SET_RETVAL(fork());
}

void
_cgo_libc_waitpid(argset_t* x) {
	pid_t pid = (pid_t)x->args[0];
	int* status = (int*)x->args[1];
	int options = (int)x->args[2];
	SET_RETVAL(waitpid(pid, status, options));
}

void
_cgo_libc_setuid(argset_t* x) {
	uid_t uid = (uid_t)x->args[0];
	SET_RETVAL(setuid(uid));
}

void
_cgo_libc_seteuid(argset_t* x) {
	uid_t euid = (uid_t)x->args[0];
	SET_RETVAL(seteuid(euid));
}

void
_cgo_libc_setgid(argset_t* x) {
	gid_t gid = (gid_t)x->args[0];
	SET_RETVAL(setgid(gid));
}

void
_cgo_libc_setegid(argset_t* x) {
	gid_t egid = (gid_t)x->args[0];
	SET_RETVAL(setegid(egid));
}

void
_cgo_libc_setreuid(argset_t* x) {
	uid_t ruid = (uid_t)x->args[0];
	uid_t euid = (uid_t)x->args[1];
	SET_RETVAL(setreuid(ruid, euid));
}

void
_cgo_libc_setregid(argset_t* x) {
	gid_t rgid = (gid_t)x->args[0];
	gid_t egid = (gid_t)x->args[1];
	SET_RETVAL(setregid(rgid, egid));
}

void
_cgo_libc_setgroups(argset_t* x) {
	size_t size = (size_t)x->args[0];
	const gid_t* list = (const gid_t*)x->args[1];
	SET_RETVAL(setgroups(size, list));
}

void
_cgo_libc_getpriority(argset_t* x) {
	int which = (int)x->args[0];
	id_t who = (id_t)x->args[1];
	int ret;
	// getpriority can return -1 as a valid value, so we must check errno.
	errno = 0;
	ret = getpriority(which, who);
	if (ret == -1 && errno != 0) {
		x->retval = (uintptr_t)errno;
	} else {
		// POSIX requires that priority is returned as `20 - ret`.
		x->retval = (uintptr_t)(20 - ret);
	}
}

void
_cgo_libc_setpriority(argset_t* x) {
	int which = (int)x->args[0];
	id_t who = (id_t)x->args[1];
	int prio = (int)x->args[2];
	SET_RETVAL(setpriority(which, who, prio));
}

void
_cgo_libc_getrlimit(argset_t* x) {
	int resource = (int)x->args[0];
	struct rlimit* rlim = (struct rlimit*)x->args[1];
	SET_RETVAL(getrlimit(resource, rlim));
}

void
_cgo_libc_setrlimit(argset_t* x) {
	int resource = (int)x->args[0];
	const struct rlimit* rlim = (const struct rlimit*)x->args[1];
	SET_RETVAL(setrlimit(resource, rlim));
}

void
_cgo_libc_getrusage(argset_t* x) {
	int who = (int)x->args[0];
	struct rusage* usage = (struct rusage*)x->args[1];
	SET_RETVAL(getrusage(who, usage));
}

void
_cgo_libc_setsid(argset_t* x) {
	SET_RETVAL(setsid());
}

void
_cgo_libc_setpgid(argset_t* x) {
	pid_t pid = (pid_t)x->args[0];
	pid_t pgid = (pid_t)x->args[1];
	SET_RETVAL(setpgid(pid, pgid));
}

void
_cgo_libc_sched_yield(argset_t* x) {
	SET_RETVAL(sched_yield());
}

void
_cgo_libc_exit(argset_t* x) {
	int status = (int)x->args[0];
	_exit(status);
}

void
_cgo_libc_execve(argset_t* x) {
	const char* pathname = (const char*)x->args[0];
	char* const* argv = (char* const*)x->args[1];
	char* const* envp = (char* const*)x->args[2];
	SET_RETVAL(execve(pathname, argv, envp));
}
// --- Signal Handling ---

void
_cgo_libc_kill(argset_t* x) {
	pid_t pid = (pid_t)x->args[0];
	int sig = (int)x->args[1];
	SET_RETVAL(kill(pid, sig));
}

void
_cgo_libc_raise(argset_t* x) {
	int sig = (int)x->args[0];
	SET_RETVAL(raise(sig));
}

void
_cgo_libc_sigaction(argset_t* x) {
	int signum = (int)x->args[0];
	const struct sigaction* act = (const struct sigaction*)x->args[1];
	struct sigaction* oldact = (struct sigaction*)x->args[2];
	SET_RETVAL(sigaction(signum, act, oldact));
}

void
_cgo_libc_sigaltstack(argset_t* x) {
	const stack_t* ss = (const stack_t*)x->args[0];
	stack_t* old_ss = (stack_t*)x->args[1];
	SET_RETVAL(sigaltstack(ss, old_ss));
}

void
_cgo_libc_sigprocmask(argset_t* x) {
	int how = (int)x->args[0];
	const sigset_t* set = (const sigset_t*)x->args[1];
	sigset_t* oldset = (sigset_t*)x->args[2];
	SET_RETVAL(sigprocmask(how, set, oldset));
}

// --- Memory Management ---

void
_cgo_libc_mmap(argset_t* x) {
	void* addr = (void*)x->args[0];
	size_t length = (size_t)x->args[1];
	int prot = (int)x->args[2];
	int flags = (int)x->args[3];
	int fd = (int)x->args[4];
	off_t offset = (off_t)x->args[5];
	x->retval = (uintptr_t)mmap(addr, length, prot, flags, fd, offset);
	// mmap() returns MAP_FAILED on error, which is (void*)-1
	if ((void*)x->retval == MAP_FAILED) {
		x->retval = (uintptr_t)errno;
	}
}

void
_cgo_libc_munmap(argset_t* x) {
	void* addr = (void*)x->args[0];
	size_t length = (size_t)x->args[1];
	SET_RETVAL(munmap(addr, length));
}

void
_cgo_libc_madvise(argset_t* x) {
	void* addr = (void*)x->args[0];
	size_t length = (size_t)x->args[1];
	int advice = (int)x->args[2];
	SET_RETVAL(madvise(addr, length, advice));
}

void
_cgo_libc_malloc(argset_t* x) {
	size_t size = (size_t)x->args[0];
	x->retval = (uintptr_t)malloc(size);
}

// --- Network Operations ---

void
_cgo_libc_socket(argset_t* x) {
	int domain = (int)x->args[0];
	int type = (int)x->args[1];
	int protocol = (int)x->args[2];
	SET_RETVAL(socket(domain, type, protocol));
}

void
_cgo_libc_socketpair(argset_t* x) {
	int domain = (int)x->args[0];
	int type = (int)x->args[1];
	int protocol = (int)x->args[2];
	int* sv = (int*)x->args[3];
	SET_RETVAL(socketpair(domain, type, protocol, sv));
}

void
_cgo_libc_bind(argset_t* x) {
	int sockfd = (int)x->args[0];
	const struct sockaddr* addr = (const struct sockaddr*)x->args[1];
	socklen_t addrlen = (socklen_t)x->args[2];
	SET_RETVAL(bind(sockfd, addr, addrlen));
}

void
_cgo_libc_connect(argset_t* x) {
	int sockfd = (int)x->args[0];
	const struct sockaddr* addr = (const struct sockaddr*)x->args[1];
	socklen_t addrlen = (socklen_t)x->args[2];
	SET_RETVAL(connect(sockfd, addr, addrlen));
}

void
_cgo_libc_listen(argset_t* x) {
	int sockfd = (int)x->args[0];
	int backlog = (int)x->args[1];
	SET_RETVAL(listen(sockfd, backlog));
}

void
_cgo_libc_accept(argset_t* x) {
	int sockfd = (int)x->args[0];
	struct sockaddr* addr = (struct sockaddr*)x->args[1];
	socklen_t* addrlen = (socklen_t*)x->args[2];
	SET_RETVAL(accept(sockfd, addr, addrlen));
}

void
_cgo_libc_shutdown(argset_t* x) {
	int sockfd = (int)x->args[0];
	int how = (int)x->args[1];
	SET_RETVAL(shutdown(sockfd, how));
}

void
_cgo_libc_sendto(argset_t* x) {
	int sockfd = (int)x->args[0];
	const void* buf = (const void*)x->args[1];
	size_t len = (size_t)x->args[2];
	int flags = (int)x->args[3];
	const struct sockaddr* dest_addr = (const struct sockaddr*)x->args[4];
	socklen_t addrlen = (socklen_t)x->args[5];
	SET_RETVAL(sendto(sockfd, buf, len, flags, dest_addr, addrlen));
}

void
_cgo_libc_sendmsg(argset_t* x) {
	int sockfd = (int)x->args[0];
	const struct msghdr* msg = (const struct msghdr*)x->args[1];
	int flags = (int)x->args[2];
	SET_RETVAL(sendmsg(sockfd, msg, flags));
}

void
_cgo_libc_recvfrom(argset_t* x) {
	int sockfd = (int)x->args[0];
	void* buf = (void*)x->args[1];
	size_t len = (size_t)x->args[2];
	int flags = (int)x->args[3];
	struct sockaddr* src_addr = (struct sockaddr*)x->args[4];
	socklen_t* addrlen = (socklen_t*)x->args[5];
	SET_RETVAL(recvfrom(sockfd, buf, len, flags, src_addr, addrlen));
}

void
_cgo_libc_recvmsg(argset_t* x) {
	int sockfd = (int)x->args[0];
	struct msghdr* msg = (struct msghdr*)x->args[1];
	int flags = (int)x->args[2];
	SET_RETVAL(recvmsg(sockfd, msg, flags));
}

void
_cgo_libc_getsockopt(argset_t* x) {
	int sockfd = (int)x->args[0];
	int level = (int)x->args[1];
	int optname = (int)x->args[2];
	void* optval = (void*)x->args[3];
	socklen_t* optlen = (socklen_t*)x->args[4];
	SET_RETVAL(getsockopt(sockfd, level, optname, optval, optlen));
}

void
_cgo_libc_setsockopt(argset_t* x) {
	int sockfd = (int)x->args[0];
	int level = (int)x->args[1];
	int optname = (int)x->args[2];
	const void* optval = (const void*)x->args[3];
	socklen_t optlen = (socklen_t)x->args[4];
	SET_RETVAL(setsockopt(sockfd, level, optname, optval, optlen));
}

void
_cgo_libc_getpeername(argset_t* x) {
	int sockfd = (int)x->args[0];
	struct sockaddr* addr = (struct sockaddr*)x->args[1];
	socklen_t* addrlen = (socklen_t*)x->args[2];
	SET_RETVAL(getpeername(sockfd, addr, addrlen));
}

void
_cgo_libc_getsockname(argset_t* x) {
	int sockfd = (int)x->args[0];
	struct sockaddr* addr = (struct sockaddr*)x->args[1];
	socklen_t* addrlen = (socklen_t*)x->args[2];
	SET_RETVAL(getsockname(sockfd, addr, addrlen));
}

// --- Epoll Operations ---

void
_cgo_libc_epoll_create1(argset_t* x) {
	int flags = (int)x->args[0];
	SET_RETVAL(epoll_create1(flags));
}

void
_cgo_libc_epoll_ctl(argset_t* x) {
	int epfd = (int)x->args[0];
	int op = (int)x->args[1];
	int fd = (int)x->args[2];
	struct epoll_event* event = (struct epoll_event*)x->args[3];
	SET_RETVAL(epoll_ctl(epfd, op, fd, event));
}

void
_cgo_libc_epoll_wait(argset_t* x) {
	int epfd = (int)x->args[0];
	struct epoll_event* events = (struct epoll_event*)x->args[1];
	int maxevents = (int)x->args[2];
	int timeout = (int)x->args[3];
	SET_RETVAL(epoll_wait(epfd, events, maxevents, timeout));
}


// --- Time Operations ---

void
_cgo_libc_gettimeofday(argset_t* x) {
	struct timeval* tv = (struct timeval*)x->args[0];
	SET_RETVAL(gettimeofday(tv, NULL)); // tz is obsolete
}

void
_cgo_libc_nanosleep(argset_t* x) {
	const struct timespec* req = (const struct timespec*)x->args[0];
	struct timespec* rem = (struct timespec*)x->args[1];
	SET_RETVAL(nanosleep(req, rem));
}

void
_cgo_libc_clock_gettime(argset_t* x) {
	clockid_t clk_id = (clockid_t)x->args[0];
	struct timespec* tp = (struct timespec*)x->args[1];
	SET_RETVAL(clock_gettime(clk_id, tp));
}

void
_cgo_libc_setitimer(argset_t* x) {
	int which = (int)x->args[0];
	const struct itimerval* new_value = (const struct itimerval*)x->args[1];
	struct itimerval* old_value = (struct itimerval*)x->args[2];
	SET_RETVAL(setitimer(which, new_value, old_value));
}

void
_cgo_libc_adjtime(argset_t* x) {
	x->retval = (uintptr_t)-1;
	x->error = ENOSYS;
}

void
_cgo_libc_usleep(argset_t* x) {
	useconds_t usec = (useconds_t)x->args[0];
	SET_RETVAL(usleep(usec));
}


// --- Pthread Operations ---

void
_cgo_libc_pthread_attr_destroy(argset_t* x) {
	pthread_attr_t* attr = (pthread_attr_t*)x->args[0];
	SET_PTHREAD_RETVAL(pthread_attr_destroy(attr));
}

void
_cgo_libc_pthread_attr_getstack(argset_t* x) {
	const pthread_attr_t* attr = (const pthread_attr_t*)x->args[0];
	void** stackaddr = (void**)x->args[1];
	size_t* stacksize = (size_t*)x->args[2];
	SET_PTHREAD_RETVAL(pthread_attr_getstack(attr, stackaddr, stacksize));
}

void
_cgo_libc_pthread_attr_init(argset_t* x) {
	pthread_attr_t* attr = (pthread_attr_t*)x->args[0];
	SET_PTHREAD_RETVAL(pthread_attr_init(attr));
}

void
_cgo_libc_pthread_attr_setdetachstate(argset_t* x) {
	pthread_attr_t* attr = (pthread_attr_t*)x->args[0];
	int detachstate = (int)x->args[1];
	SET_PTHREAD_RETVAL(pthread_attr_setdetachstate(attr, detachstate));
}

void
_cgo_libc_pthread_attr_setstack(argset_t* x) {
	pthread_attr_t* attr = (pthread_attr_t*)x->args[0];
	void* stackaddr = (void*)x->args[1];
	size_t stacksize = (size_t)x->args[2];
	SET_PTHREAD_RETVAL(pthread_attr_setstack(attr, stackaddr, stacksize));
}

void
_cgo_libc_pthread_create(argset_t* x) {
	pthread_t* thread = (pthread_t*)x->args[0];
	const pthread_attr_t* attr = (const pthread_attr_t*)x->args[1];
	void* (*start_routine)(void*) = (void* (*)(void*))x->args[2];
	void* arg = (void*)x->args[3];
	SET_PTHREAD_RETVAL(pthread_create(thread, attr, start_routine, arg));
}

void
_cgo_libc_pthread_self(argset_t* x) {
	x->retval = (uintptr_t)pthread_self();
}

void
_cgo_libc_pthread_kill(argset_t* x) {
	pthread_t thread = (pthread_t)x->args[0];
	int sig = (int)x->args[1];
	SET_PTHREAD_RETVAL(pthread_kill(thread, sig));
}


// --- Semaphore Operations ---

void
_cgo_libc_sem_init(argset_t* x) {
	sem_t* sem = (sem_t*)x->args[0];
	int pshared = (int)x->args[1];
	unsigned int value = (unsigned int)x->args[2];
	SET_RETVAL(sem_init(sem, pshared, value));
}

void
_cgo_libc_sem_post(argset_t* x) {
	sem_t* sem = (sem_t*)x->args[0];
	SET_RETVAL(sem_post(sem));
}

void
_cgo_libc_sem_timedwait(argset_t* x) {
	sem_t* sem = (sem_t*)x->args[0];
	const struct timespec* abs_timeout = (const struct timespec*)x->args[1];
	SET_RETVAL(sem_timedwait(sem, abs_timeout));
}

void
_cgo_libc_sem_wait(argset_t* x) {
	sem_t* sem = (sem_t*)x->args[0];
	SET_RETVAL(sem_wait(sem));
}


// --- System Information ---

void
_cgo_libc_uname(argset_t* x) {
	struct utsname* buf = (struct utsname*)x->args[0];
	SET_RETVAL(uname(buf));
}

void
_cgo_libc_gethostname(argset_t* x) {
	char* name = (char*)x->args[0];
	size_t len = (size_t)x->args[1];
	SET_RETVAL(gethostname(name, len));
}

void
_cgo_libc_getexecname(argset_t* x) {
	x->retval = 0;
	x->error = ENOSYS;
}

void
_cgo_libc_sysconf(argset_t* x) {
	int name = (int)x->args[0];
	long ret;
	// sysconf can return -1 as a valid value
	errno = 0;
	ret = sysconf(name);
	if (ret == -1 && errno != 0) {
		x->retval = (uintptr_t)errno;
	} else {
		x->retval = (uintptr_t)ret;
	}
}

extern char** environ;
void
_cgo_libc_environ(argset_t* x) {
	x->retval = (uintptr_t)environ;
}
