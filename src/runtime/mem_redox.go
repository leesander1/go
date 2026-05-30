// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package runtime

import "unsafe"

const _ENOMEM = 12

// Don't split the stack as this function may be invoked without a valid G,
// which prevents us from allocating more stack.
//
//go:nosplit
func sysAllocOS(n uintptr, _ string) unsafe.Pointer {
	v, err := mmap(nil, n, _PROT_READ|_PROT_WRITE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
	if err != 0 {
		return nil
	}
	return v
}

func sysUnusedOS(v unsafe.Pointer, n uintptr) {
}

func sysUsedOS(v unsafe.Pointer, n uintptr) {
}

func sysHugePageOS(v unsafe.Pointer, n uintptr) {
}

func sysNoHugePageOS(v unsafe.Pointer, n uintptr) {
}

func sysHugePageCollapseOS(v unsafe.Pointer, n uintptr) {
}

// Don't split the stack as this function may be invoked without a valid G,
// which prevents us from allocating more stack.
//
//go:nosplit
func sysFreeOS(v unsafe.Pointer, n uintptr) {
	munmap(v, n)
}

func sysFaultOS(v unsafe.Pointer, n uintptr) {
	if errno := mprotect(v, n, _PROT_NONE); errno != 0 {
		print("runtime: mprotect(", v, ", ", n, ") returned ", errno, "\n")
		throw("runtime: cannot fault pages in arena address space")
	}
}

func sysReserveOS(v unsafe.Pointer, n uintptr, _ string) unsafe.Pointer {
	p, err := mmap(v, n, _PROT_NONE, _MAP_ANON|_MAP_PRIVATE, -1, 0)
	if err != 0 {
		return nil
	}
	return p
}

func sysMapOS(v unsafe.Pointer, n uintptr, _ string) {
	if errno := mprotect(v, n, _PROT_READ|_PROT_WRITE); errno != 0 {
		if errno == _ENOMEM {
			throw("runtime: out of memory")
		}
		print("runtime: mprotect(", v, ", ", n, ") returned ", errno, "\n")
		throw("runtime: cannot map pages in arena address space")
	}
}

func needZeroAfterSysUnusedOS() bool {
	return true
}
