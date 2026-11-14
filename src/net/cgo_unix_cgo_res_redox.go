// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo && !netgo

package net

/*
#include <sys/types.h>

extern int res_search(const char *, int, int, unsigned char *, int);
*/
import "C"

type _C_struct___res_state = struct{}

func _C_res_ninit(state *_C_struct___res_state) error {
	return nil
}

func _C_res_nclose(state *_C_struct___res_state) {
}

func _C_res_nsearch(state *_C_struct___res_state, dname *_C_char, class, typ int, ans *_C_uchar, anslen int) int {
	return int(C.res_search(dname, C.int(class), C.int(typ), ans, C.int(anslen)))
}
