// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <pthread.h>
#include <errno.h>
#include <string.h>
#include <signal.h>
#include <stdlib.h>
#include "libcgo.h"
#include "libcgo_unix.h"

static void* threadentry(void*);
static void (*setg_gcc)(void*);

void
x_cgo_init(G *g, void (*setg)(void*))
{
	uintptr *pbounds;

	setg_gcc = setg;
	pbounds = (uintptr*)malloc(2 * sizeof(uintptr));
	if (pbounds == NULL) {
		fprintf(stderr, "runtime/cgo: malloc failed: %s\n", strerror(errno));
		abort();
	}
	_cgo_set_stacklo(g, pbounds);
	free(pbounds);
}

void
_cgo_sys_thread_start(ThreadStart *ts)
{
	pthread_attr_t attr;
	sigset_t ign, oset;
	pthread_t p;
	void *base;
	size_t size;
	int err;

	sigfillset(&ign);
	pthread_sigmask(SIG_SETMASK, &ign, &oset);

	pthread_attr_init(&attr);

	if (pthread_attr_getstack(&attr, &base, &size) != 0)
		perror("runtime/cgo: pthread_attr_getstack failed");
	if (size == 0) {
		ts->g->stackhi = 2 << 20;
		if (pthread_attr_setstack(&attr, NULL, ts->g->stackhi) != 0)
			perror("runtime/cgo: pthread_attr_setstack failed");
	} else {
		ts->g->stackhi = size;
	}
	// pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);

	err = _cgo_try_pthread_create(&p, &attr, threadentry, ts);

	pthread_sigmask(SIG_SETMASK, &oset, nil);


	if (err != 0) {
		perror("pthread_create failed");
	}
}

extern void crosscall1(void (*fn)(void), void (*setg_gcc)(void*), void *g);
static void*
threadentry(void *v)
{
	ThreadStart ts;

	ts = *(ThreadStart*)v;
	free(v);

	crosscall1(ts.fn, setg_gcc, (void*)ts.g);
	return nil;
}
