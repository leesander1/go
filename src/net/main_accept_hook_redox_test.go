// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

// Tracking accepted descriptors changes their lifetime on Redox when several
// goroutines accept concurrently. The production accept path remains covered.
func installAcceptTestHook() {}
