// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT _rt0_arm64_redox(SB),NOSPLIT,$-8
	JMP	_rt0_arm64(SB)

TEXT _rt0_arm64_redox_lib(SB),NOSPLIT,$0
	JMP	_rt0_arm64_lib(SB)
