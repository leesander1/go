// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !redox

package runtime_test

import (
	"internal/testenv"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestGOMAXPROCSUpdateHelper(t *testing.T) {
	const child = "GO_RUNTIME_TEST_MAXPROCS_HELPER"
	if os.Getenv(child) == "1" {
		deadline := time.Now().Add(10 * time.Second)
		for {
			if runtime.GOMAXPROCSUpdateHelperReady() {
				return
			}
			if time.Now().After(deadline) {
				t.Fatal("automatic GOMAXPROCS helper did not start")
			}
			time.Sleep(time.Millisecond)
		}
	}
	testenv.MustHaveExec(t)
	cmd := exec.Command(os.Args[0], "-test.run=^TestGOMAXPROCSUpdateHelper$", "-test.timeout=20s")
	cmd.Env = append(os.Environ(), child+"=1", "GODEBUG=updatemaxprocs=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("helper subprocess: %v\n%s", err, out)
	}
}
